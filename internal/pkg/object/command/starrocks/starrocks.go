package starrocks

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow/flight"
	"github.com/apache/arrow-go/v18/arrow/flight/flightsql"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	serviceName    = "starrocks"
	authHeaderName = "authorization"
)

var (
	handleMethod         = telemetry.NewMethod("handle", serviceName)
	createExcMethod      = telemetry.NewMethod("createExc", serviceName)
	collectResultsMethod = telemetry.NewMethod("collectResults", serviceName)
)

type clusterContext struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Database string `yaml:"database,omitempty" json:"database,omitempty"`
	UseTLS   bool   `yaml:"use_tls,omitempty" json:"use_tls,omitempty"`
}

type commandContext struct {
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

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
		creds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})
	} else {
		creds = insecure.NewCredentials()
	}

	flightClient, err := flight.NewClientWithMiddlewareCtx(ctx, endpoint, nil, nil, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to open StarRocks Flight SQL connection to %q: %v", endpoint, err)
	}

	return &flightsql.Client{Client: flightClient, Alloc: memory.DefaultAllocator}, nil

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
