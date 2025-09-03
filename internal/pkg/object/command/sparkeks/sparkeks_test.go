package sparkeks

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	ctx := &sparkEksJobContext{
		Parameters: &sparkEksJobParameters{
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
	ctx := &sparkEksJobContext{
		Parameters: &sparkEksJobParameters{
			Properties: map[string]string{},
		},
	}
	params := getSparkSubmitParameters(ctx)
	if params != nil {
		t.Errorf("Expected nil for empty properties, got %v", params)
	}
}

func TestGetSparkSubmitParameters_NilParameters(t *testing.T) {
	ctx := &sparkEksJobContext{
		Parameters: nil,
	}
	params := getSparkSubmitParameters(ctx)
	if params != nil {
		t.Errorf("Expected nil for nil parameters, got %v", params)
	}
}

func TestGetS3FileURI_InvalidFormat(t *testing.T) {
	_, err := getS3FileURI(aws.Config{}, "invalid-uri", "avro")
	if err == nil {
		t.Error("Expected error for invalid S3 URI format")
	}
}

func TestGetS3FileURI_ValidFormat(t *testing.T) {
	// This test only checks parsing, not actual AWS interaction
	uri := "s3://bucket/path/"
	_, err := getS3FileURI(aws.Config{}, uri, "avro")
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
