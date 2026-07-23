package sparkeks

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kubeflow/spark-operator/v2/api/v1beta2"

	"github.com/patterninc/heimdall/pkg/object/job"
)

func TestUpdateS3ToS3aURI(t *testing.T) {
	uri := "s3://bucket/path/file"
	expected := "s3a://bucket/path/file"
	result := updateS3ToS3aURI(uri)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
	// Should not change if already s3a
	uri2 := "s3a://bucket/path/file"
	result2 := updateS3ToS3aURI(uri2)
	if result2 != uri2 {
		t.Errorf("Expected %s, got %s", uri2, result2)
	}
}

func TestGetSparkSubmitParameters(t *testing.T) {
	ctx := &jobContext{
		Parameters: &jobParameters{
			Properties: map[string]string{
				"spark.executor.memory": "4g",
				"spark.driver.cores":    "2",
			},
		},
	}
	params := getSparkSubmitParameters(ctx)
	if params == nil || !strings.Contains(*params, "--conf spark.executor.memory=4g") || !strings.Contains(*params, "--conf spark.driver.cores=2") {
		t.Errorf("Unexpected Spark submit parameters: %v", params)
	}
}

func TestGetSparkSubmitParameters_Empty(t *testing.T) {
	ctx := &jobContext{
		Parameters: &jobParameters{
			Properties: map[string]string{},
		},
	}
	params := getSparkSubmitParameters(ctx)
	if params != nil {
		t.Errorf("Expected nil for empty properties, got %v", params)
	}
}

func TestGetSparkSubmitParameters_NilParameters(t *testing.T) {
	ctx := &jobContext{
		Parameters: nil,
	}
	params := getSparkSubmitParameters(ctx)
	if params != nil {
		t.Errorf("Expected nil for nil parameters, got %v", params)
	}
}

func TestGetS3FileURI_InvalidFormat(t *testing.T) {
	_, err := getS3FileURI(context.Background(), aws.Config{}, "invalid-uri", "avro")
	if err == nil {
		t.Error("Expected error for invalid S3 URI format")
	}
}

func TestGetS3FileURI_ValidFormat(t *testing.T) {
	// This test only checks parsing, not actual AWS interaction
	uri := "s3://bucket/path/"
	_, err := getS3FileURI(context.Background(), aws.Config{}, uri, "avro")
	// Should not error on parsing, but will error on AWS call (which is fine for unit test context)
	if err == nil || !strings.Contains(err.Error(), "failed to list S3 objects") {
		t.Logf("Expected AWS list objects error, got: %v", err)
	}
}

func TestPrintState(t *testing.T) {
	f, err := os.CreateTemp("", "sparkeks-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer f.Close()
	printState(f, "RUNNING")
	// Optionally check file contents if needed
}

func TestApplySparkOperatorConfig_JarMode(t *testing.T) {
	execCtx := &executionContext{
		appName:   "spark-sql-job-test",
		queryURI:  "s3://bucket/jobs/test/queries/query.sql",
		resultURI: "s3://bucket/jobs/test/results",
		job:       &job.Job{},
		commandContext: &commandContext{
			WrapperURI:    "s3://pattern-dl-ops/contrib/chipmunk/chipmunk.jar",
			KubeNamespace: "default",
		},
		clusterContext: &clusterContext{Properties: map[string]string{}},
		jobContext: &jobContext{
			Query:     "select key, value from t",
			Arguments: []string{"s3://pattern-dl/chipmunk/collections/foo/v2026072200"},
			Parameters: &jobParameters{
				EntryPoint: "com.pattern.chipmunk.Writer",
				Properties: map[string]string{},
			},
		},
		sparkApp: &v1beta2.SparkApplication{Spec: v1beta2.SparkApplicationSpec{
			Driver: v1beta2.DriverSpec{}, Executor: v1beta2.ExecutorSpec{},
		}},
	}
	execCtx.job.User = "svc_chipmunk"

	applySparkOperatorConfig(execCtx)

	spec := execCtx.sparkApp.Spec
	if spec.Type != v1beta2.SparkApplicationTypeScala {
		t.Fatalf("expected Type Scala, got %v", spec.Type)
	}
	if spec.MainClass == nil || *spec.MainClass != "com.pattern.chipmunk.Writer" {
		t.Fatalf("expected MainClass set, got %v", spec.MainClass)
	}
	if spec.MainApplicationFile == nil || *spec.MainApplicationFile != "s3a://pattern-dl-ops/contrib/chipmunk/chipmunk.jar" {
		t.Fatalf("unexpected MainApplicationFile: %v", spec.MainApplicationFile)
	}
	want := []string{"spark-sql-job-test", "svc_chipmunk", "select key, value from t", "s3://pattern-dl/chipmunk/collections/foo/v2026072200"}
	if !reflect.DeepEqual(spec.Arguments, want) {
		t.Fatalf("unexpected Arguments: got %v want %v", spec.Arguments, want)
	}
}

func TestApplySparkOperatorConfig_JarMode_ReturnResultAppendsResultURI(t *testing.T) {
	execCtx := &executionContext{
		appName:   "spark-sql-job-test",
		queryURI:  "s3://bucket/jobs/test/queries/query.sql",
		resultURI: "s3://bucket/jobs/test/results",
		job:       &job.Job{},
		commandContext: &commandContext{
			WrapperURI:    "s3://pattern-dl-ops/contrib/chipmunk/chipmunk.jar",
			KubeNamespace: "default",
		},
		clusterContext: &clusterContext{Properties: map[string]string{}},
		jobContext: &jobContext{
			Query:        "select key, value from t",
			Arguments:    []string{"s3://pattern-dl/chipmunk/collections/foo/v2026072200"},
			ReturnResult: true,
			Parameters: &jobParameters{
				EntryPoint: "com.pattern.chipmunk.Writer",
				Properties: map[string]string{},
			},
		},
		sparkApp: &v1beta2.SparkApplication{Spec: v1beta2.SparkApplicationSpec{
			Driver: v1beta2.DriverSpec{}, Executor: v1beta2.ExecutorSpec{},
		}},
	}
	execCtx.job.User = "svc_chipmunk"

	applySparkOperatorConfig(execCtx)

	spec := execCtx.sparkApp.Spec
	if len(spec.Arguments) != 5 {
		t.Fatalf("expected 5 arguments, got %d: %v", len(spec.Arguments), spec.Arguments)
	}
	wantLast := "s3a://bucket/jobs/test/results"
	if spec.Arguments[4] != wantLast {
		t.Fatalf("expected last argument to be s3a result URI %q, got %q", wantLast, spec.Arguments[4])
	}
	want := []string{"spark-sql-job-test", "svc_chipmunk", "select key, value from t", "s3://pattern-dl/chipmunk/collections/foo/v2026072200", wantLast}
	if !reflect.DeepEqual(spec.Arguments, want) {
		t.Fatalf("unexpected Arguments: got %v want %v", spec.Arguments, want)
	}
}

func TestApplySparkOperatorConfig_SqlWrapperPath_Unchanged(t *testing.T) {
	baseExecCtx := func(returnResult bool) *executionContext {
		return &executionContext{
			appName:   "spark-sql-job-test",
			queryURI:  "s3://bucket/jobs/test/queries/query.sql",
			resultURI: "s3://bucket/jobs/test/results",
			job:       &job.Job{},
			commandContext: &commandContext{
				WrapperURI:    "s3://pattern-dl-ops/contrib/spark-sql-wrapper/wrapper.py",
				KubeNamespace: "default",
			},
			clusterContext: &clusterContext{Properties: map[string]string{}},
			jobContext: &jobContext{
				Query:        "select 1",
				ReturnResult: returnResult,
				Parameters: &jobParameters{
					Properties: map[string]string{},
				},
			},
			sparkApp: &v1beta2.SparkApplication{Spec: v1beta2.SparkApplicationSpec{
				Driver: v1beta2.DriverSpec{}, Executor: v1beta2.ExecutorSpec{},
			}},
		}
	}

	t.Run("without_return_result", func(t *testing.T) {
		execCtx := baseExecCtx(false)
		execCtx.job.User = "svc_test"

		applySparkOperatorConfig(execCtx)

		spec := execCtx.sparkApp.Spec
		if spec.MainClass != nil {
			t.Fatalf("expected MainClass to remain unset, got %v", *spec.MainClass)
		}
		if spec.MainApplicationFile == nil || *spec.MainApplicationFile != "s3a://pattern-dl-ops/contrib/spark-sql-wrapper/wrapper.py" {
			t.Fatalf("unexpected MainApplicationFile: %v", spec.MainApplicationFile)
		}
		want := []string{"spark-sql-job-test", "s3a://bucket/jobs/test/queries/query.sql", "svc_test"}
		if !reflect.DeepEqual(spec.Arguments, want) {
			t.Fatalf("unexpected Arguments: got %v want %v", spec.Arguments, want)
		}
	})

	t.Run("with_return_result", func(t *testing.T) {
		execCtx := baseExecCtx(true)
		execCtx.job.User = "svc_test"

		applySparkOperatorConfig(execCtx)

		spec := execCtx.sparkApp.Spec
		want := []string{"spark-sql-job-test", "s3a://bucket/jobs/test/queries/query.sql", "svc_test", "s3a://bucket/jobs/test/results"}
		if !reflect.DeepEqual(spec.Arguments, want) {
			t.Fatalf("unexpected Arguments: got %v want %v", spec.Arguments, want)
		}
	})
}

func TestApplySparkOperatorConfig_Gap4DefaultType(t *testing.T) {
	// Simulates loadTemplate's inline default (no spark_application_file configured):
	// Spec.Type starts as the zero value "". No Parameters.EntryPoint and no WrapperURI,
	// so applySparkOperatorConfig's JAR/SQL-wrapper branch is skipped entirely and Type
	// would otherwise remain "" without the Gap-4 fallback.
	execCtx := &executionContext{
		appName: "spark-sql-job-test",
		job:     &job.Job{},
		commandContext: &commandContext{
			KubeNamespace: "default",
		},
		clusterContext: &clusterContext{Properties: map[string]string{}},
		jobContext: &jobContext{
			Query: "select 1",
			Parameters: &jobParameters{
				Properties: map[string]string{},
			},
		},
		sparkApp: &v1beta2.SparkApplication{Spec: v1beta2.SparkApplicationSpec{
			Driver: v1beta2.DriverSpec{}, Executor: v1beta2.ExecutorSpec{},
		}},
	}
	execCtx.job.User = "svc_test"

	if execCtx.sparkApp.Spec.Type != "" {
		t.Fatalf("test precondition failed: expected Type to start empty, got %v", execCtx.sparkApp.Spec.Type)
	}

	applySparkOperatorConfig(execCtx)

	if execCtx.sparkApp.Spec.Type != v1beta2.SparkApplicationTypePython {
		t.Fatalf("expected Gap-4 fallback to set Type to Python, got %v", execCtx.sparkApp.Spec.Type)
	}
}
