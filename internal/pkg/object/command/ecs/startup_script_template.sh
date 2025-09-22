#!/bin/bash
set -e
echo 'Starting file downloads to {{.DownloadDir}}...'
{{if .CreateDirs}}mkdir -p {{.DownloadDir}}{{end}}

if ! command -v aws &> /dev/null; then
  echo 'Installing AWS CLI...'
  if apk update && apk add aws-cli; then
    echo 'AWS CLI installed successfully'
  else
    echo 'ERROR: Failed to install AWS CLI'
    exit 1
  fi
fi

{{range .Downloads}}
# Download: {{.Source}}
mkdir -p $(dirname {{.Destination}})
{{if .IsS3}}
echo "Downloading from S3: {{.Source}}"
if aws s3 cp '{{.Source}}' '{{.Destination}}' --cli-read-timeout {{$.Timeout}} --cli-connect-timeout {{$.Timeout}}; then
  echo "Successfully downloaded: {{.Source}}"
  
  if [ -f '{{.Destination}}' ] && [ -s '{{.Destination}}' ]; then
    echo "File verification passed: {{.Destination}}"
    file_size=$(stat -c%s '{{.Destination}}' 2>/dev/null || echo "unknown")
    echo "File size: $file_size bytes"
  else
    echo "ERROR: Downloaded file is empty or missing: {{.Destination}}"
    exit 1
  fi
else
  echo "ERROR: Failed to download from S3: {{.Source}}"
  exit 1
fi
{{end}}
{{end}}
echo 'All files downloaded successfully'
echo 'Starting main application...'
# Execute the original command
if [ -n "$ORIGINAL_COMMAND" ]; then
  echo "Executing: $ORIGINAL_COMMAND"
  exec $ORIGINAL_COMMAND
else
  echo "No original command found, executing: $@"
  exec "$@"
fi
