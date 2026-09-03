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
			execCtx := newTestExecutionContext(tt.wrapperURI, &jobContext{}, "alice", "result_uri")
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
	execCtx := newTestExecutionContext("s3://bucket/app.jar", jobCtx, "alice", "result_uri")

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
	expectedArgs := []string{"spark-sql-job-test", "s3a://bucket/query.sql", "alice", "result_uri"}
	if !reflect.DeepEqual(spec.Arguments, expectedArgs) {
		t.Errorf("spec.Arguments = %v, want %v", spec.Arguments, expectedArgs)
	}
}

func TestNewSQLWrapperEntrypointStrategy_Apply(t *testing.T) {
	jobCtx := &jobContext{ReturnResult: true}
	execCtx := newTestExecutionContext("s3://bucket/wrapper.py", jobCtx, "bob", "result_uri")

	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Fatalf("apply() returned unexpected error: %v", err)
	}

	if spec.MainClass != nil {
		t.Errorf("spec.MainClass = %v, want nil for sql wrapper", spec.MainClass)
	}
	expectedArgs := []string{"spark-sql-job-test", "s3a://bucket/query.sql", "bob", "result_uri"}
	if !reflect.DeepEqual(spec.Arguments, expectedArgs) {
		t.Errorf("spec.Arguments = %v, want %v", spec.Arguments, expectedArgs)
	}
}

// Caller-supplied arguments extend the managed prefix rather than replacing it: user keeps slot 2
// and extras keep a stable index regardless of return_result.
func TestBuildArguments_ExtrasAppendedAtStableIndex(t *testing.T) {
	tests := []struct {
		name         string
		returnResult bool
		expected     []string
	}{
		{
			name:         "with result",
			returnResult: true,
			expected:     []string{"app", "s3a://bucket/query.sql", "alice", "s3a://bucket/result", "s3://bucket/prefix/", "--flag"},
		},
		{
			name:         "without result holds the slot",
			returnResult: false,
			expected:     []string{"app", "s3a://bucket/query.sql", "alice", "", "s3://bucket/prefix/", "--flag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extra := []string{"s3://bucket/prefix/", "--flag"}
			got := buildArguments(extra, "app", "s3a://bucket/query.sql", "alice", "s3a://bucket/result", tt.returnResult)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("buildArguments() = %v, want %v", got, tt.expected)
			}
			if got[2] != "alice" {
				t.Errorf("user = %q, want alice at index 2", got[2])
			}
			if got[4] != extra[0] {
				t.Errorf("first extra = %q, want %q at index 4", got[4], extra[0])
			}
		})
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
			execCtx := newTestExecutionContext("s3://bucket/app.jar", tt.jobCtx, "alice", "result_uri")
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
	execCtx := newTestExecutionContext("s3://bucket/wrapper.py", &jobContext{}, "alice", "result_uri")
	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Errorf("apply() returned unexpected error: %v", err)
	}
}

func newTestExecutionContext(wrapperURI string, jobCtx *jobContext, user string, resultURI string) *executionContext {
	return &executionContext{
		job: &job.Job{
			Object: object.Object{User: user},
		},
		commandContext: &commandContext{WrapperURI: wrapperURI},
		jobContext:     jobCtx,
		appName:        "spark-sql-job-test",
		s3aQueryURI:    "s3a://bucket/query.sql",
		s3aResultURI:   resultURI,
	}
}

func newTestBundleExecutionContext(bundleURI string, jobCtx *jobContext, user string, resultURI string) *executionContext {
	execCtx := newTestExecutionContext("s3://bucket/pyspark-job-wrapper.py", jobCtx, user, resultURI)
	execCtx.commandContext.BundleURI = bundleURI
	return execCtx
}

func TestNewEntrypointStrategy_BundleURISelectsPySparkStrategy(t *testing.T) {
	execCtx := newTestBundleExecutionContext("s3://bucket/pyspark/buybox-predictor/", &jobContext{}, "alice", "result_uri")
	strategy := newEntrypointStrategy(execCtx)
	if got := reflect.TypeOf(strategy); got != reflect.TypeOf(pysparkEntrypointStrategy{}) {
		t.Errorf("newEntrypointStrategy() type = %v, want pysparkEntrypointStrategy", got)
	}
}

func TestNewEntrypointStrategy_ScriptURISelectsPySparkStrategy(t *testing.T) {
	jobCtx := &jobContext{
		Parameters: &jobParameters{
			ScriptURI: "s3://bucket/pyspark/scripts/alice/train.py",
		},
	}
	execCtx := newTestBundleExecutionContext("s3://bucket/pyspark/", jobCtx, "alice", "result_uri")
	strategy := newEntrypointStrategy(execCtx)
	if got := reflect.TypeOf(strategy); got != reflect.TypeOf(pysparkEntrypointStrategy{}) {
		t.Errorf("newEntrypointStrategy() type = %v, want pysparkEntrypointStrategy", got)
	}
}

func TestPySparkBundleEntrypointStrategy_Apply(t *testing.T) {
	jobCtx := &jobContext{
		ReturnResult: true,
		Arguments:    []string{"--date=2026-08-31"},
		Parameters: &jobParameters{
			EntryPoint:    "src/prediction_flow.py",
			BundleVersion: "2f3a91c",
		},
	}
	execCtx := newTestBundleExecutionContext("s3://bucket/pyspark/buybox-predictor/", jobCtx, "alice", "result_uri")

	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Fatalf("apply() returned unexpected error: %v", err)
	}

	expectedArgs := []string{
		"spark-sql-job-test", "s3a://bucket/query.sql", "alice", "result_uri",
		"s3a://bucket/pyspark/buybox-predictor/2f3a91c.zip", "src/prediction_flow.py",
		"--date=2026-08-31",
	}
	if !reflect.DeepEqual(spec.Arguments, expectedArgs) {
		t.Errorf("spec.Arguments = %v, want %v", spec.Arguments, expectedArgs)
	}
}

func TestPySparkScriptEntrypointStrategy_Apply(t *testing.T) {
	jobCtx := &jobContext{
		ReturnResult: true,
		Arguments:    []string{"--date=2026-08-31"},
		Parameters: &jobParameters{
			ScriptURI: "s3://bucket/pyspark/scripts/alice/train.py",
		},
	}
	execCtx := newTestBundleExecutionContext("s3://bucket/pyspark/", jobCtx, "alice", "result_uri")

	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Fatalf("apply() returned unexpected error: %v", err)
	}

	expectedArgs := []string{
		"spark-sql-job-test", "s3a://bucket/query.sql", "alice", "result_uri",
		"s3a://bucket/pyspark/scripts/alice/train.py", "",
		"--date=2026-08-31",
	}
	if !reflect.DeepEqual(spec.Arguments, expectedArgs) {
		t.Errorf("spec.Arguments = %v, want %v", spec.Arguments, expectedArgs)
	}
}

func TestValidateScriptURI(t *testing.T) {
	prefix := "s3://bucket/pyspark/"
	tests := []struct {
		name      string
		scriptURI string
		wantErr   error
	}{
		{name: "valid", scriptURI: "s3://bucket/pyspark/scripts/alice/train.py"},
		{name: "valid s3a", scriptURI: "s3a://bucket/pyspark/scripts/alice/train.py"},
		{name: "wrong prefix", scriptURI: "s3://other-bucket/pyspark/train.py", wantErr: ErrInvalidScriptURI},
		{name: "not python", scriptURI: "s3://bucket/pyspark/scripts/alice/train.sh", wantErr: ErrInvalidScriptURI},
		{name: "traversal", scriptURI: "s3://bucket/pyspark/../evil.py", wantErr: ErrInvalidScriptURI},
		{name: "not s3", scriptURI: "https://bucket/pyspark/train.py", wantErr: ErrInvalidScriptURI},
		{name: "empty", scriptURI: "", wantErr: ErrInvalidScriptURI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScriptURI(tt.scriptURI, prefix)
			if tt.wantErr == nil && err != nil {
				t.Fatalf("validateScriptURI() error = %v, want nil", err)
			}
			if tt.wantErr != nil && err != tt.wantErr {
				t.Errorf("validateScriptURI() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPySparkBundleEntrypointStrategy_Apply_MissingFields(t *testing.T) {
	tests := []struct {
		name    string
		jobCtx  *jobContext
		wantErr error
	}{
		{name: "nil parameters", jobCtx: &jobContext{}, wantErr: ErrMissingBundleVersion},
		{name: "missing bundle_version", jobCtx: &jobContext{Parameters: &jobParameters{EntryPoint: "src/main.py"}}, wantErr: ErrMissingBundleVersion},
		{name: "missing entry_point", jobCtx: &jobContext{Parameters: &jobParameters{BundleVersion: "abc123"}}, wantErr: ErrMissingBundleEntry},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCtx := newTestBundleExecutionContext("s3://bucket/pyspark/proj/", tt.jobCtx, "alice", "result_uri")
			strategy := newEntrypointStrategy(execCtx)
			spec := &v1beta2.SparkApplicationSpec{}
			if err := strategy.apply(spec); err != tt.wantErr {
				t.Errorf("apply() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPySparkBundleEntrypointStrategy_Apply_ProjectSlashVersion(t *testing.T) {
	jobCtx := &jobContext{
		Parameters: &jobParameters{
			EntryPoint:    "src/update_prediction_data.py",
			BundleVersion: "buybox-predictor/2f3a91c",
		},
	}
	execCtx := newTestBundleExecutionContext("s3://bucket/pyspark/", jobCtx, "alice", "")

	strategy := newEntrypointStrategy(execCtx)
	spec := &v1beta2.SparkApplicationSpec{}
	if err := strategy.apply(spec); err != nil {
		t.Fatalf("apply() returned unexpected error: %v", err)
	}
	if spec.Arguments[4] != "s3a://bucket/pyspark/buybox-predictor/2f3a91c.zip" {
		t.Errorf("bundle uri = %q", spec.Arguments[4])
	}
}
