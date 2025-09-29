package heimdall

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	metricsProxyScheme = `http`
	localHost          = `127.0.0.1`
)

func metricsRouteHandler(w http.ResponseWriter, r *http.Request) {
	metricsPortAddress := os.Getenv(`PROMETHEUS_ADDRESS`)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: metricsProxyScheme,
		Host:   localHost + metricsPortAddress,
	})
	proxy.ServeHTTP(w, r)
}

func repeatOverDimensionsWithParam[T any](fn func(T, ...string), parameter T, command, cluster string) {
	fn(parameter)                   // overall
	fn(parameter, cluster)          // per cluster
	fn(parameter, command)          // per command
	fn(parameter, cluster, command) // per cluster and command
}

func repeatOverDimensions(fn func(...string), command, cluster string) {
	fn()                 // overall
	fn(cluster)          // per cluster
	fn(command)          // per command
	fn(cluster, command) // per cluster and command
}
