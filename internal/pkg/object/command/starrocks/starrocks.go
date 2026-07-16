package starrocks

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow/flight"
	"github.com/apache/arrow-go/v18/arrow/flight/flightsql"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
