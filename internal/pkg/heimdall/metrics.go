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
