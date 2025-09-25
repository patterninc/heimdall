package heimdall

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/patterninc/heimdall/internal/pkg/auth"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/internal/pkg/janitor"
	"github.com/patterninc/heimdall/internal/pkg/pool"
	"github.com/patterninc/heimdall/internal/pkg/server"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/command"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultJobsDirectory    = `/tmp/heimdall`
	defaultArchiveDirectory = defaultJobsDirectory + `/archive`
	defaultResultDirectory  = defaultJobsDirectory + `/result`
	formatErrUnknownPlugin  = "unknown plugin: %s"
	defaultAPIPrefix        = `/api/v1`
	methodPOST              = `POST`
	methodPUT               = `PUT`
	methodGET               = `GET`
	webUIProxyScheme        = `http`
	webUIProxyHost          = `127.0.0.1:3000`
	webUIUsernameHeader     = `X-Heimdall-User`
	webUIVersionHeader      = `X-Heimdall-Version`
	defaultUsername         = `user@heimdall`
)

type Heimdall struct {
	Server           *server.Server       `yaml:"server,omitempty" json:"server,omitempty"`
	Commands         command.Commands     `yaml:"commands,omitempty" json:"commands,omitempty"`
	Clusters         cluster.Clusters     `yaml:"clusters,omitempty" json:"clusters,omitempty"`
	JobsDirectory    string               `yaml:"jobs_directory,omitempty" json:"jobs_directory,omitempty"`
	ArchiveDirectory string               `yaml:"archive_directory,omitempty" json:"archive_directory,omitempty"`
	ResultDirectory  string               `yaml:"result_directory,omitempty" json:"result_directory,omitempty"`
	PluginsDirectory string               `yaml:"plugin_directory,omitempty" json:"plugin_directory,omitempty"`
	Database         *database.Database   `yaml:"database,omitempty" json:"database,omitempty"`
	Pool             *pool.Pool[*job.Job] `yaml:"pool,omitempty" json:"pool,omitempty"`
	Auth             *auth.Auth           `yaml:"auth,omitempty" json:"auth,omitempty"`
	Janitor          *janitor.Janitor     `yaml:"janitor,omitempty" json:"janitor,omitempty"`
	Version          string               `yaml:"-" json:"-"`
	agentName        string
	commandHandlers  map[string]plugin.Handler
}

func (h *Heimdall) Init() error {

	// set jobs directory if not set
	if h.JobsDirectory == `` {
		h.JobsDirectory = defaultJobsDirectory
	}

	// set archive directory if not set
	if h.ArchiveDirectory == `` {
		h.ArchiveDirectory = defaultArchiveDirectory
	}

	// set result directory if not set
	if h.ResultDirectory == `` {
		h.ResultDirectory = defaultResultDirectory
	}

	// let's set the agent name
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	h.agentName = fmt.Sprintf("%s-%d", strings.ToLower(hostname), time.Now().UnixMicro())

	// let's load all the plugins
	plugins, err := h.loadPlugins()
	if err != nil {
		return err
	}

	h.commandHandlers = make(map[string]plugin.Handler)

	// process commands / add default values if missing, write commands to db
	for _, c := range h.Commands {

		// set defaults for missing properties
		if err := c.Init(); err != nil {
			return err
		}

		// set command handlers
		pluginNew, found := plugins[c.Plugin]
		if !found {
			return fmt.Errorf(formatErrUnknownPlugin, c.Plugin)
		}

		handler, err := pluginNew(c.Context)
		if err != nil {
			return err
		}
		h.commandHandlers[c.ID] = handler

		// let's record command in the database
		if err := h.commandUpsert(c); err != nil {
			return err
		}

	}

	// process commands / add default values if missing, write commands to db
	for _, c := range h.Clusters {

		// set defaults for missing properties
		if err := c.Init(); err != nil {
			return err
		}

		// let's record command in the database
		if err := h.clusterUpsert(c); err != nil {
			return err
		}

	}

	// start janitor
	if err := h.Janitor.Start(h.Database); err != nil {
		return err
	}

	// let's start the agent
	return h.Pool.Start(h.runAsyncJob, h.getAsyncJobs)

}

func (h *Heimdall) Start() error {
	// set routes
	router := mux.NewRouter()

	// let's set auth middleware so we could use it to call a plugin
	router.Use(h.auth)

	// set API routes
	apiRouter := router.PathPrefix(defaultAPIPrefix).Subrouter()

	// job(s) endpoints...
	apiRouter.Methods(methodGET).PathPrefix(`/job/statuses`).HandlerFunc(payloadHandler(h.getJobStatuses))
	apiRouter.Methods(methodGET).PathPrefix(`/job/{id}/status`).HandlerFunc(payloadHandler(h.getJobStatus))
	apiRouter.Methods(methodGET).PathPrefix(`/job/{id}/{file}`).HandlerFunc(h.getJobFile)
	apiRouter.Methods(methodGET).PathPrefix(`/job/{id}`).HandlerFunc(payloadHandler(h.getJob))
	apiRouter.Methods(methodGET).PathPrefix(`/jobs`).HandlerFunc(payloadHandler(h.getJobs))
	apiRouter.Methods(methodPOST).PathPrefix(`/job`).HandlerFunc(payloadHandler(h.submitJob))
	apiRouter.Methods(methodGET).PathPrefix(`/command/statuses`).HandlerFunc(payloadHandler(h.getCommandStatuses))
	apiRouter.Methods(methodGET).PathPrefix(`/command/{id}/status`).HandlerFunc(payloadHandler(h.getCommandStatus))
	apiRouter.Methods(methodPUT).PathPrefix(`/command/{id}/status`).HandlerFunc(payloadHandler(h.updateCommandStatus))
	apiRouter.Methods(methodPUT).PathPrefix(`/command/{id}`).HandlerFunc(payloadHandler(h.submitCommand))
	apiRouter.Methods(methodGET).PathPrefix(`/command/{id}`).HandlerFunc(payloadHandler(h.getCommand))
	apiRouter.Methods(methodGET).PathPrefix(`/commands`).HandlerFunc(payloadHandler(h.getCommands))
	apiRouter.Methods(methodGET).PathPrefix(`/cluster/statuses`).HandlerFunc(payloadHandler(h.getClusterStatuses))
	apiRouter.Methods(methodGET).PathPrefix(`/cluster/{id}/status`).HandlerFunc(payloadHandler(h.getClusterStatus))
	apiRouter.Methods(methodPUT).PathPrefix(`/cluster/{id}/status`).HandlerFunc(payloadHandler(h.updateClusterStatus))
	apiRouter.Methods(methodPUT).PathPrefix(`/cluster/{id}`).HandlerFunc(payloadHandler(h.submitCluster))
	apiRouter.Methods(methodGET).PathPrefix(`/cluster/{id}`).HandlerFunc(payloadHandler(h.getCluster))
	apiRouter.Methods(methodGET).PathPrefix(`/clusters`).HandlerFunc(payloadHandler(h.getClusters))

	// metrics endpoint - proxy to metrics service
	router.Path(`/metrics`).HandlerFunc(metricsRouteHandler)

	// catch all for APIs
	apiRouter.PathPrefix(`/`).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, fmt.Errorf("unknown endpoint: %s %s", r.Method, r.URL.Path), nil)
	})

	// pass everything else to nextjs Web UI
	router.PathPrefix(`/`).Handler(http.StripPrefix(`/`, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: webUIProxyScheme,
				Host:   webUIProxyHost,
			})

			// pass the username to web UI
			username := defaultUsername
			if u := getUsername(r); u != `` {
				username = u
			}

			r.Header.Set(webUIUsernameHeader, username)
			r.Header.Set(webUIVersionHeader, h.Version)

			proxy.ServeHTTP(w, r)

		}),
	))

	return h.Server.Start(router)

}
