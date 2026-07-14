package starrocks

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/apache/arrow-go/v18/arrow/flight"
	"github.com/apache/arrow-go/v18/arrow/flight/flightsql"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

type commandContext struct {
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

type clusterContext struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Database string `yaml:"database,omitempty" json:"database,omitempty"`
	UseTLS   bool   `yaml:"use_tls,omitempty" json:"use_tls,omitempty"`
}

type jobContext struct {
	Query        string `yaml:"query" json:"query"`
	ReturnResult bool   `yaml:"return_result,omitempty" json:"return_result,omitempty"`
	client       *flightsql.Client
	authToken    string
	endpoint     string
	useTLS       bool
}

const (
	serviceName    = "starrocks"
	authHeaderName = "authorization"
)

var (
	handleMethod         = telemetry.NewMethod("handle", serviceName)
	createExcMethod      = telemetry.NewMethod("createExc", serviceName)
	collectResultsMethod = telemetry.NewMethod("collectResults", serviceName)
)

// New creates a new starrocks plugin handler
func New(commandCtx *heimdallContext.Context) (plugin.Handler, error) {
	t := &commandContext{}

	if commandCtx != nil {
		if err := commandCtx.Unmarshal(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Execute implements the plugin.Handler interface
func (cmd *commandContext) Execute(ctx context.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	handleMethod.CountRequest()
	defer handleMethod.RecordLatency(time.Now())

	jobContext, err := cmd.createJobContext(ctx, j, c)
	if err != nil {
		handleMethod.LogAndCountError(err, "create_job_context")
		return err
	}
	defer jobContext.close()

	res, err := jobContext.execute(ctx)
	if err != nil {
		handleMethod.LogAndCountError(err, "execute")
		return err
	}
	j.Result = res
	j.Status = status.Succeeded

	handleMethod.CountSuccess()
	return nil
}

func (cmd *commandContext) createJobContext(ctx context.Context, j *job.Job, c *cluster.Cluster) (*jobContext, error) {

	clusterCtx := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterCtx); err != nil {
			createExcMethod.CountError("unmarshal_cluster_context")
			return nil, fmt.Errorf("failed to unmarshal cluster context: %v", err)
		}
	}

	jobCtx := &jobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobCtx); err != nil {
			createExcMethod.CountError("unmarshal_job_context")
			return nil, fmt.Errorf("failed to unmarshal job context: %v", err)
		}
	}

	client, token, err := connect(ctx, clusterCtx.Endpoint, clusterCtx.UseTLS, cmd.Username, cmd.Password)
	if err != nil {
		createExcMethod.CountError("open_connection")
		return nil, err
	}

	jobCtx.client = client
	jobCtx.authToken = token
	jobCtx.endpoint = clusterCtx.Endpoint
	jobCtx.useTLS = clusterCtx.UseTLS

	return jobCtx, nil

}

// connect performs the basic-auth handshake and returns the client plus the
// bearer token required on every subsequent RPC.
func connect(ctx context.Context, endpoint string, useTLS bool, username, password string) (*flightsql.Client, string, error) {

	flightClient, err := dialFlightClient(ctx, endpoint, useTLS)
	if err != nil {
		return nil, "", err
	}

	authCtx, err := flightClient.Client.AuthenticateBasicToken(ctx, username, password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to authenticate with StarRocks Flight SQL server: %v", err)
	}

	token := ``
	if md, ok := metadata.FromOutgoingContext(authCtx); ok {
		if values := md.Get(authHeaderName); len(values) > 0 {
			token = values[0]
		}
	}

	return flightClient, token, nil

}

// dialFlightClient opens a plain gRPC connection to a Flight SQL endpoint
// (FE or BE); BE connections just need the bearer token in the context, not a fresh handshake.
func dialFlightClient(ctx context.Context, endpoint string, useTLS bool) (*flightsql.Client, error) {

	var creds credentials.TransportCredentials
	if useTLS {
		creds = credentials.NewTLS(&tls.Config{})
	} else {
		creds = insecure.NewCredentials()
	}

	flightClient, err := flight.NewClientWithMiddlewareCtx(ctx, endpoint, nil, nil, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to open StarRocks Flight SQL connection to %q: %v", endpoint, err)
	}

	return &flightsql.Client{Client: flightClient, Alloc: memory.DefaultAllocator}, nil

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

// backendClientPool lazily dials one Flight SQL client per unique BE
// location, reusing the FE client when an endpoint has no distinct BE address.
type backendClientPool struct {
	feClient       *flightsql.Client
	dialedEndpoint string
	useTLS         bool

	mu      sync.Mutex
	clients map[string]*flightsql.Client
}

func newBackendClientPool(feClient *flightsql.Client, dialedEndpoint string, useTLS bool) *backendClientPool {
	return &backendClientPool{
		feClient:       feClient,
		dialedEndpoint: dialedEndpoint,
		useTLS:         useTLS,
		clients:        make(map[string]*flightsql.Client),
	}
}

// clientFor returns the FE client, or a cached/dialed direct BE connection if the endpoint has a distinct location.
func (p *backendClientPool) clientFor(ctx context.Context, endpoint *flight.FlightEndpoint) (*flightsql.Client, error) {

	for _, loc := range endpoint.Location {
		if loc.Uri == flight.LocationReuseConnection {
			continue
		}
		host := locationHost(loc.Uri)
		if host == `` {
			return nil, fmt.Errorf("failed to parse flight endpoint location %q", loc.Uri)
		}
		if host == p.dialedEndpoint {
			continue
		}
		return p.get(ctx, host)
	}

	return p.feClient, nil

}

func (p *backendClientPool) get(ctx context.Context, host string) (*flightsql.Client, error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if c, ok := p.clients[host]; ok {
		return c, nil
	}

	c, err := dialFlightClient(ctx, host, p.useTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to StarRocks backend %q: %v", host, err)
	}

	p.clients[host] = c
	return c, nil

}

func (p *backendClientPool) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.clients {
		c.Client.Close()
	}
}

func fetchEndpoint(ctx context.Context, beClients *backendClientPool, endpoint *flight.FlightEndpoint) ([][]any, error) {

	client, err := beClients.clientFor(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	reader, err := client.DoGet(ctx, endpoint.Ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch result stream: %v", err)
	}
	defer reader.Release()

	rows := make([][]any, 0, 128)
	for reader.Next() {
		rows = append(rows, recordToRows(reader.Record())...)
	}

	if err := reader.Err(); err != nil {
		return nil, fmt.Errorf("error reading result stream: %v", err)
	}

	return rows, nil

}

// locationHost extracts host:port from a Flight location URI (e.g. "grpc+tcp://host:9408").
func locationHost(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ``
	}
	return u.Host
}

// HealthCheck implements the plugin.HealthChecker interface
func (cmd *commandContext) HealthCheck(ctx context.Context, c *cluster.Cluster) error {

	clusterCtx := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterCtx); err != nil {
			return err
		}
	}

	client, token, err := connect(ctx, clusterCtx.Endpoint, clusterCtx.UseTLS, cmd.Username, cmd.Password)
	if err != nil {
		return err
	}
	defer client.Client.Close()

	ctx = metadata.AppendToOutgoingContext(ctx, authHeaderName, token)

	info, err := client.Execute(ctx, "SELECT 1")
	if err != nil {
		return err
	}

	beClients := newBackendClientPool(client, clusterCtx.Endpoint, clusterCtx.UseTLS)
	defer beClients.closeAll()

	for _, endpoint := range info.Endpoint {
		if _, err := fetchEndpoint(ctx, beClients, endpoint); err != nil {
			return err
		}
	}

	return nil

}

func (cmd *commandContext) Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error {
	// No cleanup needed. StarRocks queries should always be synchronous.
	return nil
}
