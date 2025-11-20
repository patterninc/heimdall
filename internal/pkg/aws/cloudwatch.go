package aws

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func PullLogs(ctx context.Context, writer *os.File, logGroup, logStream string, chunkSize int, memoryLimit int64) error {

	// initialize AWS session
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// create a new cloudwatch logs client
	client := cloudwatchlogs.NewFromConfig(cfg)

	// Write log stream header to differentiate streams
	header := fmt.Sprintf("=== %s ===\n\n", logStream)
	writer.WriteString(header)

	totalMemoryUsed := int64(len(header))

	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		StartFromHead: aws.Bool(true),
		Limit:         aws.Int32(int32(chunkSize)),
	}

	var previousToken *string

	// Loop through log events
	for {
		output, err := client.GetLogEvents(ctx, input)
		if err != nil {
			return err
		}

		// Write in chunks to the writer
		for _, event := range output.Events {
			// Format log with timestamp and message only
			timestamp := time.UnixMilli(*event.Timestamp).UTC().Format("2006-01-02 15:04:05")
			message := aws.ToString(event.Message)

			formattedLog := fmt.Sprintf("[%s] %s\n", timestamp, message)
			writer.WriteString(formattedLog)

			// Track memory usage: formatted log length
			totalMemoryUsed += int64(len(formattedLog))
		}

		// Stop if the memory limit is reached
		if totalMemoryUsed >= memoryLimit {
			writer.WriteString("Memory limit reached. Truncating logs.\n")
			break
		}

		// If the next token is the same as the previous token, we have reached the end of the logs
		if previousToken != nil && *previousToken == *output.NextForwardToken {
			break
		}

		previousToken = output.NextForwardToken
		input.NextToken = output.NextForwardToken

		// Small pause to avoid rate limiting
		time.Sleep(200 * time.Millisecond)
	}

	return nil

}
