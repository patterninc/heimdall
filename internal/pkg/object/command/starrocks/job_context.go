package starrocks

import (
	"context"
	"fmt"

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

func (j *jobContext) execute(ctx context.Context) (*result.Result, error) {

	ctx = j.withAuth(ctx)

	if !j.ReturnResult {
		n, err := j.client.ExecuteUpdate(ctx, j.Query)
		if err != nil {
			return nil, fmt.Errorf("failed to execute statement: %v", err)
		}
		return result.FromMessage(fmt.Sprintf("%d row(s) affected", n))
	}

	info, err := j.client.Execute(ctx, j.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return collectResults(ctx, j.client, info, j.endpoint, j.useTLS)

}

// collectResults fans ticket redemption out across goroutines, connecting
// directly to each BE instead of funneling every row back through the FE.
func collectResults(ctx context.Context, feClient *flightsql.Client, info *flight.FlightInfo, dialedEndpoint string, useTLS bool) (*result.Result, error) {

	schema, err := flight.DeserializeSchema(info.Schema, memory.DefaultAllocator)
	if err != nil {
		collectResultsMethod.CountError("deserialize_schema")
		return nil, fmt.Errorf("failed to parse result schema: %v", err)
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
		return nil, err
	}

	out := &result.Result{
		Columns: columnsFromSchema(schema),
		Data:    make([][]any, 0, 128),
	}
	for _, rows := range rowsByEndpoint {
		out.Data = append(out.Data, rows...)
	}

	return out, nil

}
