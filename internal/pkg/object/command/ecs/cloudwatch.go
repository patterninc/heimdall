package ecs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

const (
	maxLogChunkSize   = 200                // Process 200 log entries at a time
	maxLogMemoryBytes = 1024 * 1024 * 1024 // 1GB safety limit
)

func pullCloudWatchLogs(client *cloudwatchlogs.Client, writer *os.File, logGroup, logStream string) error {

	// Write log stream header to differentiate streams
	header := fmt.Sprintf("=== %s ===\n\n", logStream)
	writer.WriteString(header)

	totalMemoryUsed := int64(len(header))

	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		StartFromHead: aws.Bool(true),
		Limit:         aws.Int32(maxLogChunkSize),
	}

	var previousToken *string

	// Loop through log events
	for {
		output, err := client.GetLogEvents(context.Background(), input)
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

		// Stop if: memory limit reached or same token returned (end of stream)
		if totalMemoryUsed >= maxLogMemoryBytes ||
			(previousToken != nil && *previousToken == *output.NextForwardToken) {
			break
		}

		previousToken = output.NextForwardToken
		input.NextToken = output.NextForwardToken

		// Small pause to avoid rate limiting
		time.Sleep(200 * time.Millisecond)
	}

	return nil

}
