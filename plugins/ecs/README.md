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

## 🖥️ Cluster Configuration

ECS clusters must specify Fargate configuration, IAM roles, and network settings:

```yaml
- name: fargate-cluster
  status: active
  version: 0.0.1
  description: ECS Fargate cluster for batch jobs
  context:
    max_cpu: 4096
    max_memory: 8192
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
    "cpu": 256,
    "memory": 512,
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

## 📊 Job Logs

The ECS plugin supports CloudWatch log retrieval from your containers after job completion. Configure the `awslogs` log driver in your task definition to enable log collection:

```json
{
  "logConfiguration": {
    "logDriver": "awslogs",
    "options": {
      "awslogs-group": "/ecs/batch-job",
      "awslogs-region": "us-west-2", 
      "awslogs-stream-prefix": "ecs"
    }
  }
}
```

Logs are formatted with timestamps and stream context:

```
=== ecs/main/abc123def456789 ===

[2023-11-11 13:45:23] Starting application...
[2023-11-11 13:45:24] Processing data batch
[2023-11-11 13:45:25] Task completed successfully
```

**Log Routing:**
- **Successful jobs**: Logs from a successful task → **stdout**
- **Failed jobs**: Logs from the failed task (max retries) → **stderr**  
- **Timeout jobs**: Logs from an incomplete task → **stderr**

**Note**: When multiple tasks run in parallel, only logs from a single representative task are retrieved after the job completes. For successful jobs, this is the last completed task. For failed jobs, this is the task that reached maximum retries. 

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
    cpu: 256
    memory: 512
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