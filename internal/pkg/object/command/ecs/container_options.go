package ecs

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// ContainerOption is a function that modifies container overrides
type ContainerOption func(*executionContext) error

// ApplyContainerOptions applies a list of container options to container overrides
func ApplyContainerOptions(execCtx *executionContext, options ...ContainerOption) error {
	for _, option := range options {
		if err := option(execCtx); err != nil {
			return err
		}
	}
	return nil
}

// WithFileUploadScript wraps container commands with a script that downloads files from S3 to the container
func WithFileUploadScript(fileUploads []FileUpload, localDir string) ContainerOption {
	return func(execCtx *executionContext) error {
		if len(fileUploads) == 0 {
			return nil
		}

		for i := range execCtx.ContainerOverrides {
			override := &execCtx.ContainerOverrides[i]

			// Get the original command from override, or from task definition if not overridden
			var originalCommand []string
			if len(override.Command) > 0 {
				originalCommand = override.Command
			} else {
				// Get command from task definition
				for _, container := range execCtx.TaskDefinitionWrapper.TaskDefinition.ContainerDefinitions {
					if aws.ToString(container.Name) == aws.ToString(override.Name) {
						originalCommand = container.Command
						break
					}
				}
			}

			if len(originalCommand) == 0 {
				originalCommand = []string{}
			}

			// Generate the download wrapper script
			wrapperScript := generateDownloadWrapperScript(fileUploads, localDir, originalCommand)

			// Replace container command with wrapper script
			override.Command = []string{"sh", "-c", wrapperScript}
		}

		return nil
	}
}

const downloadScriptTemplate = `
set -e
apk update && apk add aws-cli
mkdir -p {{LOCAL_DIR}}
for s3_path in {{S3_PATHS}};do aws s3 cp "$s3_path" "{{LOCAL_DIR}}/$(basename "$s3_path")" 2>&1;done
exec {{CMD}}`

// generateDownloadWrapperScript generates a minimal bash script that downloads files from S3 and executes the original command
func generateDownloadWrapperScript(fileUploads []FileUpload, localDir string, originalCommand []string) string {
	// Build S3 paths list for the for loop
	var s3Paths []string
	for _, upload := range fileUploads {
		s3Paths = append(s3Paths, fmt.Sprintf(`"%s"`, upload.S3Destination))
	}
	s3PathsList := strings.Join(s3Paths, " ")

	// Build command string
	cmdStr := "wait"
	if len(originalCommand) > 0 {
		escapedCmd := make([]string, len(originalCommand))
		for i, arg := range originalCommand {
			escapedCmd[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", "'\\''"))
		}
		cmdStr = strings.Join(escapedCmd, " ")
	}

	// Fill template
	script := strings.ReplaceAll(downloadScriptTemplate, "{{LOCAL_DIR}}", localDir)
	script = strings.ReplaceAll(script, "{{S3_PATHS}}", s3PathsList)
	script = strings.ReplaceAll(script, "{{CMD}}", cmdStr)
	return script
}

// Future container options can be added here as needed
// Example:
// func WithEnvironmentVariables(envVars map[string]string) ContainerOption { ... }
// func WithHealthCheck(config HealthCheckConfig) ContainerOption { ... }
// func WithResourceLimits(limits ResourceLimits) ContainerOption { ... }
