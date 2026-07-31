package sparkeks

import (
	"reflect"
	"testing"

	"github.com/kubeflow/spark-operator/v2/api/v1beta2"

	"github.com/patterninc/heimdall/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/job"
)

func TestResolveJarApplicationType(t *testing.T) {
	tests := []struct {
		name            string
		applicationType string
		expected        v1beta2.SparkApplicationType
	}{
		{name: "scala lowercase", applicationType: "scala", expected: v1beta2.SparkApplicationTypeScala},
		{name: "java lowercase", applicationType: "java", expected: v1beta2.SparkApplicationTypeJava},
		{name: "java mixed case", applicationType: "Java", expected: v1beta2.SparkApplicationTypeJava},
		{name: "scala uppercase", applicationType: "SCALA", expected: v1beta2.SparkApplicationTypeScala},
		{name: "java with surrounding whitespace", applicationType: "  java  ", expected: v1beta2.SparkApplicationTypeJava},
		{name: "empty defaults to scala", applicationType: "", expected: defaultJarApplicationType},
		{name: "whitespace only defaults to scala", applicationType: "   ", expected: defaultJarApplicationType},
		{name: "unrecognized defaults to scala", applicationType: "python", expected: defaultJarApplicationType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveJarApplicationType(tt.applicationType); got != tt.expected {
				t.Errorf("resolveJarApplicationType(%q) = %v, want %v", tt.applicationType, got, tt.expected)
			}
		})
	}
}
func TestNewEntrypointStrategy_Selection(t *testing.T) {
	tests := []struct {
		name         string
		wrapperURI   string
		expectedType reflect.Type
	}{
		{name: "jar selects jar strategy", wrapperURI: "s3://bucket/app.jar", expectedType: reflect.TypeOf(jarEntrypointStrategy{})},
		{name: "py selects sql wrapper strategy", wrapperURI: "s3://bucket/wrapper.py", expectedType: reflect.TypeOf(sqlWrapperEntrypointStrategy{})},
		{name: "jar uppercase extension", wrapperURI: "s3://bucket/APP.JAR", expectedType: reflect.TypeOf(jarEntrypointStrategy{})},
		{name: "py mixed case selects sql strategy", wrapperURI: "s3://bucket/Wrapper.Py", expectedType: reflect.TypeOf(sqlWrapperEntrypointStrategy{})},
		{name: "unknown extension defaults to sql wrapper", wrapperURI: "s3://bucket/app.sh", expectedType: reflect.TypeOf(sqlWrapperEntrypointStrategy{})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCtx := newTestExecutionContext(tt.wrapperURI, &jobContext{}, "alice")
			strategy := newEntrypointStrategy(execCtx)
			if got := reflect.TypeOf(strategy); got != tt.expectedType {
				t.Errorf("newEntrypointStrategy(%q) type = %v, want %v", tt.wrapperURI, got, tt.expectedType)
			}
		})
	}
}

func TestNewJarEntrypointStrategy_Apply(t *testing.T) {
	jobCtx := &jobContext{
		ReturnResult: true,
		Parameters: &jobParameters{
			EntryPoint:      "com.example.Main",
			ApplicationType: "java",
		},
	}
	execCtx := newTestExecutionContext("s3://bucket/app.jar", jobCtx, "alice")

	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Fatalf("apply() returned unexpected error: %v", err)
	}

	if spec.Type != v1beta2.SparkApplicationTypeJava {
		t.Errorf("spec.Type = %v, want %v", spec.Type, v1beta2.SparkApplicationTypeJava)
	}
	if spec.MainClass == nil || *spec.MainClass != "com.example.Main" {
		t.Errorf("spec.MainClass = %v, want com.example.Main", spec.MainClass)
	}
	expectedArgs := []string{"spark-sql-job-test", "s3a://bucket/query.sql", "alice", "s3a://bucket/result"}
	if !reflect.DeepEqual(spec.Arguments, expectedArgs) {
		t.Errorf("spec.Arguments = %v, want %v", spec.Arguments, expectedArgs)
	}
}

func TestNewSQLWrapperEntrypointStrategy_Apply(t *testing.T) {
	jobCtx := &jobContext{ReturnResult: false}
	execCtx := newTestExecutionContext("s3://bucket/wrapper.py", jobCtx, "bob")

	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Fatalf("apply() returned unexpected error: %v", err)
	}

	if spec.MainClass != nil {
		t.Errorf("spec.MainClass = %v, want nil for sql wrapper", spec.MainClass)
	}
	expectedArgs := []string{"spark-sql-job-test", "s3a://bucket/query.sql", "bob"}
	if !reflect.DeepEqual(spec.Arguments, expectedArgs) {
		t.Errorf("spec.Arguments = %v, want %v", spec.Arguments, expectedArgs)
	}
}

func TestJarEntrypointStrategy_Apply_MissingEntryPoint(t *testing.T) {
	tests := []struct {
		name   string
		jobCtx *jobContext
	}{
		{name: "nil parameters", jobCtx: &jobContext{}},
		{name: "empty entry_point", jobCtx: &jobContext{Parameters: &jobParameters{EntryPoint: ""}}},
		{name: "whitespace entry_point", jobCtx: &jobContext{Parameters: &jobParameters{EntryPoint: "   "}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCtx := newTestExecutionContext("s3://bucket/app.jar", tt.jobCtx, "alice")
			strategy := newEntrypointStrategy(execCtx)
			spec := &v1beta2.SparkApplicationSpec{}
			err := strategy.apply(spec)
			if err != ErrMissingEntryPoint {
				t.Errorf("apply() error = %v, want %v", err, ErrMissingEntryPoint)
			}
			if spec.MainClass != nil {
				t.Errorf("spec.MainClass = %v, want nil when entry_point is missing", spec.MainClass)
			}
		})
	}
}

func TestSQLWrapperEntrypointStrategy_Apply_NoEntryPointRequired(t *testing.T) {
	execCtx := newTestExecutionContext("s3://bucket/wrapper.py", &jobContext{}, "alice")
	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Errorf("apply() returned unexpected error: %v", err)
	}
}

func newTestExecutionContext(wrapperURI string, jobCtx *jobContext, user string) *executionContext {
	return &executionContext{
		job: &job.Job{
			Object: object.Object{User: user},
		},
		commandContext: &commandContext{WrapperURI: wrapperURI},
		jobContext:     jobCtx,
		appName:        "spark-sql-job-test",
		s3aQueryURI:    "s3a://bucket/query.sql",
		s3aResultURI:   "s3a://bucket/result",
	}
}
