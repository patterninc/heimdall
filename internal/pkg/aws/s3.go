package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	rxS3Path = regexp.MustCompile(`^s3://([^/]+)/(.*)$`)
)

// WriteToS3 writes a file to S3, providing the same interface as os.WriteFile function
func WriteToS3(ctx context.Context, name string, data []byte, _ os.FileMode) error {

	bucket, key, err := parseS3Path(name)
	if err != nil {
		return err
	}

	// upload file
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// Create an S3 client
	svc := s3.NewFromConfig(awsConfig)

	// Create a PutObject request to upload the file
	if _, err := svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	}); err != nil {
		return err
	}

	return nil

}

func ReadFromS3(ctx context.Context, name string) ([]byte, error) {

	bucket, key, err := parseS3Path(name)
	if err != nil {
		return nil, err
	}

	// upload file
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Create an S3 client
	svc := s3.NewFromConfig(awsConfig)

	// Create a PutObject request to upload the file
	getObjectOutput, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	if getObjectOutput != nil && getObjectOutput.Body != nil {
		defer getObjectOutput.Body.Close()
		return io.ReadAll(getObjectOutput.Body)
	}

	return nil, nil

}

func parseS3Path(name string) (string, string, error) {

	// let's parse name
	matches := rxS3Path.FindStringSubmatch(name)

	if len(matches) != 3 {
		return ``, ``, fmt.Errorf("unexpected s3 path: %s", name)
	}

	return matches[1], matches[2], nil

}
