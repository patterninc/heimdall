package starrocks

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/apache/arrow-go/v18/arrow/flight"
	"github.com/apache/arrow-go/v18/arrow/flight/flightsql"
)

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

// fetchEndpoint redeems a single BE endpoint's ticket and returns the raw Flight
// reader so the caller can stream record batches directly instead of draining
// the whole stream into memory here first.
func fetchEndpoint(ctx context.Context, beClients *backendClientPool, endpoint *flight.FlightEndpoint) (*flight.Reader, error) {

	client, err := beClients.clientFor(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	reader, err := client.DoGet(ctx, endpoint.Ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch result stream: %v", err)
	}

	return reader, nil

}

// locationHost extracts host:port from a Flight location URI (e.g. "grpc+tcp://host:9408").
func locationHost(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ``
	}
	return u.Host
}
