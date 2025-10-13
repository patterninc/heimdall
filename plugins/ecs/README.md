# 🚀 ECS Fargate Plugin

The **ECS Fargate Plugin** enables Heimdall to run AWS ECS tasks on Fargate clusters. It builds task definitions dynamically from templates and runs short-lived batch jobs with robust failure tracking and retry mechanisms.

---

## 🧩 Plugin Overview

* **Plugin Name:** `ecs`
* **Execution Mode:** Asynchronous
* **Use Case:** Running containerized batch jobs on AWS ECS Fargate with automatic retries and failure management

---

## ⚙️ Defining an ECS Command

An ECS command requires a task definition template and configuration for running tasks:

```yaml
- name: ecs-batch-job-0.0.1
  status: active
  plugin: ecs
  version: 0.0.1
  description: Run batch jobs on ECS Fargate with retry logic
  context:
    task_definition_template: /path/to/task-definition.json
    task_count: 3
    container_overrides:
      - name: main
        command: ["/usr/local/bin/my-script.sh", "arg1", "arg2"]
        environment:
          - name: ENV_VAR
            value: "value"
    polling_interval: "30s"  # 30 seconds
    timeout: "1h"            # 1 hour
    max_fail_count: 3        # max retries per task
  tags:
    - type:ecs
  cluster_tags:
    - type:fargate
```

🔸 This defines an ECS command that:
- Uses a task definition template from a local JSON file
- Runs 3 tasks in parallel
- Overrides the command and environment variables for the "main" container
- Polls for completion every 30 seconds with a 1-hour timeout
- Retries failed tasks up to 3 times before giving up

**Duration Format:**
- Use human-readable duration strings: `"30s"`, `"1h"`, `"45m"`, `"2h15m30s"`
- Examples: `"15s"`, `"1h"`, `"30m"`, `"2h30m"`, `"3600s"`

---

## 📤 File Uploads to S3

The ECS plugin supports file uploads to containers (via S3):

### Direct String Upload

Upload file content passed as strings before tasks execution starts:

```yaml
file_uploads:
  - data: "File content here as a string"
    s3_destination: "s3://my-bucket/results/output.txt"
  - data: '{"result": "success", "count": 42}'
    s3_destination: "s3://my-bucket/results/data.json"
```

**Features:**
- File content passed as string in configuration
- Uploaded before ECS tasks starts
- Uses Heimdall's IAM role credentials
- 30-second timeout per upload
- Destination must be full S3 path with filename

**Use Case:** Small metadata files, status reports, configuration files


```yaml
- name: ecs-job-with-string-uploads
  status: active
  plugin: ecs
  version: 0.0.1
  description: Run ECS tasks and upload results to S3
  context:
    task_definition_template: task.json
    task_count: 2
    
    file_uploads:
      - data: "Processing completed at 2025-10-13"
        s3_destination: "s3://results-bucket/status.txt"
      - data: '{"tasks": 2, "status": "completed"}'
        s3_destination: "s3://results-bucket/metadata.json"
```

### API usage

```json
{
  "name": "run-batch-with-upload",
  "version": "0.0.1",
  "command_criteria": ["type:ecs"],
  "cluster_criteria": ["type:fargate"],
  "context": {
    "task_count": 1,
    "file_uploads": [
      {
        "data": "Task execution results: SUCCESS",
        "destination": "s3://my-bucket/jobs/job-123/results.txt"
      },
      {
        "data": "{\"job_id\": \"123\", \"timestamp\": \"2025-10-13T10:00:00Z\", \"status\": \"completed\"}",
        "destination": "s3://my-bucket/jobs/job-123/metadata.json"
      }
    ]
  }
}
```

### S3 URI Format

File destinations must follow the S3 URI format:
- **Format**: `s3://bucket-name/path/to/filename.ext`
- **Example**: `s3://my-results/2025/10/output.json`
- The path must include the complete object key with filename

### IAM Permissions

Ensure your ECS task role has S3 upload permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::your-bucket/*"
    }
  ]
}
```

---

## 🖥️ Cluster Configuration

ECS clusters must specify Fargate configuration, IAM roles, and network settings:

```yaml
- name: fargate-cluster
  status: active
  version: 0.0.1
  description: ECS Fargate cluster for batch jobs
  context:
    cpu: 256
    memory: 512
    max_task_count: 10
    execution_role_arn: arn:aws:iam::123456789012:role/ecsTaskExecutionRole
    task_role_arn: arn:aws:iam::123456789012:role/ecsTaskRole
    cluster_name: my-fargate-cluster
    launch_type: FARGATE
    vpc_config:
      subnets:
        - subnet-12345678
        - subnet-87654321
      security_groups:
        - sg-12345678
  tags:
    - type:fargate
    - data:prod
```

---

## 🚀 Submitting an ECS Job

A typical ECS job includes task configuration and optional overrides:

```json
{
  "name": "run-batch-job",
  "version": "0.0.1",
  "command_criteria": ["type:ecs"],
  "cluster_criteria": ["type:fargate"],
  "context": {
    "task_count": 2,
    "container_overrides": [
      {
        "name": "main",
        "command": ["/usr/local/bin/process-data.sh"],
        "environment": [
          {
            "name": "DATA_SOURCE",
            "value": "s3://my-bucket/data/"
          },
          {
            "name": "BATCH_SIZE",
            "value": "1000"
          }
        ]
      }
    ]
  }
}
```

🔹 The job will:
- Run 2 tasks in parallel
- Override the command to run `process-data.sh`
- Set the `DATA_SOURCE` and `BATCH_SIZE` environment variables
- **Automatically add `TASK_NAME` environment variable** to each container (e.g., `heimdall-job-{job-id}-{task-number}`)

---

## 🔄 Failure Tracking & Retry Logic

The plugin implements robust failure tracking with the following features:

### Automatic Retries
- Failed tasks are automatically retried up to `max_fail_count` times
- Each retry creates a new task with a unique ARN
- Failed ARNs are tracked in the job results

### Failure Scenarios
1. **Task Failure**: Non-zero exit code from essential containers
2. **Timeout**: Job exceeds the configured timeout period
3. **Max Retries**: Task fails more than `max_fail_count` times

### Failure Response
- When a single task exceeds `max_fail_count`, all other running tasks are stopped
- When timeout occurs, all remaining tasks are stopped

---

## 📦 Task Definition Template

The task definition template should be a complete ECS task definition JSON file:

```json
{
  "family": "batch-job",
  "requiresCompatibilities": ["FARGATE"],
  "networkMode": "awsvpc",
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::123456789012:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "main",
      "image": "alpine:latest",
      "command": ["echo", "Hello World"],
      "essential": true,
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/batch-job",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "environment": [
        {
          "name": "ENVIRONMENT",
          "value": "production"
        },
        {
          "name": "LOG_LEVEL",
          "value": "info"
        }
      ]
    }
  ]
}
```

---

## 📊 Job Results

The plugin returns comprehensive task execution results including:

```json
{
  "columns": [
    {"name": "task_arn", "type": "string"},
    {"name": "duration", "type": "float"},
    {"name": "retries", "type": "int"},
    {"name": "failed_arns", "type": "string"}
  ],
  "data": [
    ["arn:aws:ecs:us-west-2:123456789012:task/abc123", 45.2, 0, ""],
    ["arn:aws:ecs:us-west-2:123456789012:task/def456", 43.8, 2, "arn:aws:ecs:us-west-2:123456789012:task/old1,arn:aws:ecs:us-west-2:123456789012:task/old2"]
  ]
}
```

**Result Fields:**
- `task_arn`: The final task ARN (successful or last retry)
- `duration`: Total execution time in seconds
- `retries`: Number of retries attempted
- `failed_arns`: Comma-separated list of failed task ARNs

---

## 🔍 Task Monitoring

The plugin tracks tasks using the `startedBy` tag with format `heimdall-job-{job_id}-{task_number}`. This allows for:

- Status polling to check task completion
- Retrieving task results and exit codes
- Monitoring execution times
- Tracking retry attempts and failed ARNs

---

## 🌍 Environment Variables

### Automatic Variables
The plugin automatically adds the following environment variable to all containers:
- `TASK_NAME`: The unique task identifier (e.g., `heimdall-job-abc123-0`) with the final digit being the task index

### Custom Variables
You can add custom environment variables through:
1. **Task Definition Template**: Base environment variables
2. **Container Overrides**: Container Overrides can exist in both command context and job context. Job context will override command context, and command context will override any ENVs in the task definition template.

**Example:**
```yaml
container_overrides:
  - name: "main"
    environment:
      - name: "CUSTOM_VAR"
        value: "custom_value"
      - name: "ANOTHER_VAR"
        value: "another_value"
```

**Final Environment Variables:**
- `ENVIRONMENT=production` (from template)
- `LOG_LEVEL=info` (from template)
- `CUSTOM_VAR=custom_value` (from override)
- `ANOTHER_VAR=another_value` (from override)
- `TASK_NAME=heimdall-job-abc123-0` (automatically added)

---

## 🧠 Best Practices

* **Task Definition Templates**: Use complete task definitions with proper logging configuration
* **Resource Limits**: Set appropriate CPU and memory limits based on workload requirements
* **IAM Roles**: Ensure execution and task roles have minimal necessary permissions
* **Network Configuration**: Use private subnets when possible, public subnets only if internet access is required
* **Container Overrides**: Use overrides for job-specific parameters rather than modifying templates
* **Retry Configuration**: Set appropriate `max_fail_count` based on your application's reliability
* **Timeout Settings**: Configure timeouts based on expected job duration using human-readable format
* **Environment Variables**: Use `TASK_NAME` for task-specific logic in your containers

---

## 🔐 Security Considerations

* **IAM Roles**: Use least-privilege IAM roles for task execution
* **Network Security**: Configure security groups to allow only necessary traffic
* **Container Images**: Use trusted base images and scan for vulnerabilities
* **Environment Variables**: Avoid passing sensitive data as environment variables; use AWS Secrets Manager instead
* **Task Isolation**: Each task runs with its own `TASK_NAME` for proper isolation

---

## 📝 Complete Example

### Task Definition Template (`task-definition.json`)
```json
{
  "family": "batch-job-template",
  "requiresCompatibilities": ["FARGATE"],
  "networkMode": "awsvpc",
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::123456789012:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "main",
      "image": "alpine:latest",
      "command": ["sh", "-c", "echo \"Task name: $TASK_NAME\" && exit 0"],
      "essential": true,
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/batch-jobs",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "environment": [
        {
          "name": "ENVIRONMENT",
          "value": "production"
        },
        {
          "name": "LOG_LEVEL",
          "value": "info"
        }
      ]
    }
  ],
  "ephemeralStorage": {
    "sizeInGiB": 21
  },
  "runtimePlatform": {
    "cpuArchitecture": "X86_64",
    "operatingSystemFamily": "LINUX"
  }
}
```

### Command Configuration
```yaml
- name: ecs-batch-job-0.0.1
  status: active
  plugin: ecs
  version: 0.0.1
  description: Run batch jobs on ECS Fargate with retry logic
  context:
    task_definition_template: task-definition.json
    task_count: 2
    container_overrides:
      - name: main
        environment:
          - name: CUSTOM_VAR
            value: custom_value
    polling_interval: "30s"  # 30 seconds
    timeout: "1h"            # 1 hour
    max_fail_count: 3
  tags:
    - type:ecs
  cluster_tags:
    - type:fargate
```

### Job Submission
```json
{
  "name": "run-batch-job",
  "version": "0.0.1",
  "command_criteria": ["type:ecs"],
  "cluster_criteria": ["type:fargate"],
  "context": {
    "task_count": 2,
    "container_overrides": [
      {
        "name": "main",
        "environment": [
          {
            "name": "DATA_SOURCE",
            "value": "s3://my-bucket/data/"
          }
        ]
      }
    ]
  }
}
```

This configuration will run 2 tasks, each with the environment variables:
- `ENVIRONMENT=production`
- `LOG_LEVEL=info`
- `DATA_SOURCE=s3://my-bucket/data/`
- `TASK_NAME=heimdall-job-{job-id}-{task-number}` (automatically added) 
