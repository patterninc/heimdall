package heimdall

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/command"
	"github.com/patterninc/heimdall/pkg/object/job"
	jobStatus "github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

const openAPISpecPath = `../../../api/openapi.yaml`

var openAPIMethods = []string{`get`, `put`, `post`, `delete`, `patch`, `head`, `options`, `trace`}

type openAPISchema struct {
	Properties map[string]any `yaml:"properties"`
	Enum       []string       `yaml:"enum"`
}

type openAPISpec struct {
	Servers []struct {
		URL string `yaml:"url"`
	} `yaml:"servers"`
	Paths      map[string]map[string]any `yaml:"paths"`
	Components struct {
		Schemas map[string]openAPISchema `yaml:"schemas"`
	} `yaml:"components"`
}

func loadOpenAPISpec(t *testing.T) *openAPISpec {
	t.Helper()
	data, err := os.ReadFile(openAPISpecPath)
	if err != nil {
		t.Fatalf("read %s: %v", openAPISpecPath, err)
	}
	spec := &openAPISpec{}
	if err := yaml.Unmarshal(data, spec); err != nil {
		t.Fatalf("parse %s: %v", openAPISpecPath, err)
	}
	return spec
}

func specOperations(spec *openAPISpec) []string {
	var ops []string
	for path, item := range spec.Paths {
		for key := range item {
			if slices.Contains(openAPIMethods, key) {
				ops = append(ops, strings.ToUpper(key)+` `+path)
			}
		}
	}
	sort.Strings(ops)
	return ops
}

func TestOpenAPISpecServerIsAPIPrefix(t *testing.T) {
	spec := loadOpenAPISpec(t)
	if len(spec.Servers) != 1 || spec.Servers[0].URL != defaultAPIPrefix {
		t.Fatalf("spec servers = %+v, want one server with url %s", spec.Servers, defaultAPIPrefix)
	}
}

func TestOpenAPISpecListsEveryRoute(t *testing.T) {
	spec := loadOpenAPISpec(t)

	var registered []string
	for _, route := range (&Heimdall{}).apiRoutes() {
		registered = append(registered, route.method+` `+route.path)
	}
	sort.Strings(registered)

	documented := specOperations(spec)

	for _, op := range registered {
		if !slices.Contains(documented, op) {
			t.Errorf("route %s is registered but not in %s", op, openAPISpecPath)
		}
	}
	for _, op := range documented {
		if !slices.Contains(registered, op) {
			t.Errorf("operation %s is in %s but no route is registered for it", op, openAPISpecPath)
		}
	}
}

// The API router matches path prefixes in registration order, so check that each documented path reaches its own route.
func TestOpenAPISpecPathsReachTheirRoute(t *testing.T) {
	spec := loadOpenAPISpec(t)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix(defaultAPIPrefix).Subrouter()
	for _, route := range (&Heimdall{}).apiRoutes() {
		matched := route.method + ` ` + route.path
		apiRouter.Methods(route.method).PathPrefix(route.path).HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(matched))
		})
	}

	pathValues := strings.NewReplacer(`{id}`, `some-id`, `{file}`, `stdout`)
	for _, op := range specOperations(spec) {
		method, path, _ := strings.Cut(op, ` `)
		req := httptest.NewRequest(method, defaultAPIPrefix+pathValues.Replace(path), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if got := rec.Body.String(); got != op {
			t.Errorf("%s %s reaches route %q, want %q", method, req.URL.Path, got, op)
		}
	}
}

func jsonFieldNames(t reflect.Type) []string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var names []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(`json`)
		name, _, _ := strings.Cut(tag, `,`)
		if field.Anonymous && name == `` {
			names = append(names, jsonFieldNames(field.Type)...)
			continue
		}
		if !field.IsExported() || name == `-` {
			continue
		}
		if name == `` {
			name = field.Name
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func without(names []string, drop ...string) []string {
	return slices.DeleteFunc(slices.Clone(names), func(n string) bool { return slices.Contains(drop, n) })
}

func schemaProperties(t *testing.T, spec *openAPISpec, name string) []string {
	t.Helper()
	schema, found := spec.Components.Schemas[name]
	if !found {
		t.Fatalf("schema %s is not in %s", name, openAPISpecPath)
	}
	var props []string
	for p := range schema.Properties {
		props = append(props, p)
	}
	sort.Strings(props)
	return props
}

func TestOpenAPISchemasMatchGoTypes(t *testing.T) {
	spec := loadOpenAPISpec(t)

	jobFields := jsonFieldNames(reflect.TypeFor[job.Job]())
	commandFields := jsonFieldNames(reflect.TypeFor[command.Command]())
	clusterFields := jsonFieldNames(reflect.TypeFor[cluster.Cluster]())

	// Each schema must list exactly the JSON fields of its Go type, less the fields that the handler never sets.
	exact := map[string][]string{
		`Job`:               without(jobFields, `result`),
		`SubmittedJob`:      jobFields,
		`Command`:           without(commandFields, `allowed_callers`),
		`Cluster`:           without(clusterFields, `health_check`, `rbacs`, `attributes`),
		`JobsPage`:          jsonFieldNames(reflect.TypeFor[resultset]()),
		`JobAttribute`:      jsonFieldNames(reflect.TypeFor[job.Attribute]()),
		`Result`:            jsonFieldNames(reflect.TypeFor[result.Result]()),
		`ResultColumn`:      jsonFieldNames(reflect.TypeFor[column.Column]()),
		`HealthCheckResult`: jsonFieldNames(reflect.TypeFor[healthCheckResult]()),
		`HealthChecks`:      jsonFieldNames(reflect.TypeFor[healthChecksResponse]()),
		`Error`:             {idKey, errorKey},
	}
	for name, want := range exact {
		sort.Strings(want)
		if got := schemaProperties(t, spec, name); !slices.Equal(got, want) {
			t.Errorf("schema %s has properties %v, Go type has %v", name, got, want)
		}
	}

	// Request bodies and partial responses must only use fields that the Go type decodes or encodes.
	subset := map[string][]string{
		`JobInput`:               jobFields,
		`CommandInput`:           commandFields,
		`ClusterInput`:           clusterFields,
		`JobStatusSnapshot`:      jobFields,
		`JobCancelResponse`:      jobFields,
		`ResourceStatusSnapshot`: commandFields,
		`StatusUpdate`:           commandFields,
		`CommandsPage`:           jsonFieldNames(reflect.TypeFor[resultset]()),
		`ClustersPage`:           jsonFieldNames(reflect.TypeFor[resultset]()),
	}
	for name, allowed := range subset {
		for _, prop := range schemaProperties(t, spec, name) {
			if !slices.Contains(allowed, prop) {
				t.Errorf("schema %s has property %s, which the Go type does not have", name, prop)
			}
		}
	}
}

func marshalStatuses[T json.Marshaler](t *testing.T, values []T) []string {
	t.Helper()
	var names []string
	for _, v := range values {
		data, err := v.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal status %v: %v", v, err)
		}
		var name string
		if err := json.Unmarshal(data, &name); err != nil {
			t.Fatalf("unmarshal status %s: %v", data, err)
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func TestOpenAPIStatusEnumsMatchGo(t *testing.T) {
	spec := loadOpenAPISpec(t)

	jobStatuses := []jobStatus.Status{
		jobStatus.New, jobStatus.Accepted, jobStatus.Running, jobStatus.Failed,
		jobStatus.Killed, jobStatus.Succeeded, jobStatus.Canceling, jobStatus.Canceled,
	}
	var jobMarshalers []*jobStatus.Status
	for i := range jobStatuses {
		jobMarshalers = append(jobMarshalers, &jobStatuses[i])
	}

	resourceStatuses := []status.Status{status.Active, status.Inactive, status.Deleted}
	var resourceMarshalers []*status.Status
	for i := range resourceStatuses {
		resourceMarshalers = append(resourceMarshalers, &resourceStatuses[i])
	}

	want := map[string][]string{
		`JobStatus`:      marshalStatuses(t, jobMarshalers),
		`ResourceStatus`: marshalStatuses(t, resourceMarshalers),
	}
	for name, values := range want {
		got := slices.Clone(spec.Components.Schemas[name].Enum)
		sort.Strings(got)
		if !slices.Equal(got, values) {
			t.Errorf("schema %s has enum %v, Go has %v", name, got, values)
		}
	}
}
