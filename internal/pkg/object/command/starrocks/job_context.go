package starrocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/apache/arrow-go/v18/arrow/flight"
	"github.com/apache/arrow-go/v18/arrow/flight/flightsql"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"

	"github.com/patterninc/heimdall/pkg/result"
)

type jobContext struct {
	Query        string `yaml:"query" json:"query"`
	ReturnResult bool   `yaml:"return_result,omitempty" json:"return_result,omitempty"`
	client       *flightsql.Client
	authToken    string
	endpoint     string
	useTLS       bool
}

func (j *jobContext) withAuth(ctx context.Context) context.Context {
	if j.authToken == `` {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, authHeaderName, j.authToken)
}

func (j *jobContext) close() {
	if j.client != nil {
		j.client.Client.Close()
	}
}

func (j *jobContext) execute(ctx context.Context, w io.Writer) error {

	ctx = j.withAuth(ctx)

	if !j.ReturnResult {
		n, err := j.client.ExecuteUpdate(ctx, j.Query)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %v", err)
		}
		msg, err := result.FromMessage(fmt.Sprintf("%d row(s) affected", n))
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(msg)
	}

	info, err := j.client.Execute(ctx, j.Query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return collectResults(ctx, w, j.client, info, j.endpoint, j.useTLS)

}

// collectResults fans ticket redemption out across goroutines, connecting
// directly to each BE instead of funneling every row back through the FE.
func collectResults(ctx context.Context, w io.Writer, feClient *flightsql.Client, info *flight.FlightInfo, dialedEndpoint string, useTLS bool) error {

	schema, err := flight.DeserializeSchema(info.Schema, memory.DefaultAllocator)
	if err != nil {
		collectResultsMethod.CountError("deserialize_schema")
		return fmt.Errorf("failed to parse result schema: %v", err)
	}

	beClients := newBackendClientPool(feClient, dialedEndpoint, useTLS)
	defer beClients.closeAll()

	// endpoints redeem independently, so fan DoGet out across goroutines
	// instead of draining them one at a time.
	rowsByEndpoint := make([][][]any, len(info.Endpoint))

	g, gctx := errgroup.WithContext(ctx)
	for i, endpoint := range info.Endpoint {
		g.Go(func() error {
			rows, err := fetchEndpoint(gctx, beClients, endpoint)
			if err != nil {
				return err
			}
			rowsByEndpoint[i] = rows
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		collectResultsMethod.CountError("do_get")
		return err
	}

	rw := result.NewRowWriter(w)
	if err := rw.WriteColumns(columnsFromSchema(schema)); err != nil {
		return err
	}
	for i, rows := range rowsByEndpoint {
		for _, row := range rows {
			if err := rw.WriteRow(row); err != nil {
				return err
			}
		}
		rowsByEndpoint[i] = nil // release this endpoint's rows as soon as they're written, instead of holding every endpoint alive until this whole slice goes out of scope
	}

	return rw.Close()

}
