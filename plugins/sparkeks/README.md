# Spark EKS Plugin

The Spark EKS plugin enables submitting Spark jobs to an AWS EKS cluster using the Kubeflow Spark Operator. It supports IAM role assumption, S3 integration for queries/results/logs, and custom Spark application templates.

## Features

- Submit Spark jobs to EKS via Kubeflow Spark Operator
- Upload SQL queries to S3 and fetch results from S3
- Collect Spark pod logs and store them in S3
- Supports custom SparkApplication YAML templates
- IAM role assumption for cross-account EKS access
- Configurable Spark job resources and properties

## Configuration

### Cluster Context

```json
{
  "name": "my-eks-cluster",
  "context": {
    "role_arn": "arn:aws:iam::ACCOUNT:role/EKSAccessRole",
    "region": "us-west-2",
    "image": "myrepo/spark:latest",
    "spark_application_file": "/path/to/spark-application.yaml",
    "properties": {
      "spark.hadoop.fs.s3a.access.key": "...",
      "spark.hadoop.fs.s3a.secret.key": "..."
    }
  }
}
```

### Job Context

```json
{
  "query": "SELECT * FROM my_table",
  "properties": {
    "spark.driver.memory": "2g",
    "spark.executor.instances": "2"
  },
  "return_result": true
}
```

### Command Context

```json
{
  "queries_uri": "s3://mybucket/queries",
  "results_uri": "s3://mybucket/results",
  "logs_uri": "s3://mybucket/logs",
  "wrapper_uri": "s3://mybucket/wrapper.py",
  "properties": {
    "spark.some.config": "value"
  },
  "kube_namespace": "default"
}
```

## Usage

Submit a job using the API:

```json
{
  "name": "run-spark-query",
  "version": "0.0.1",
  "command_criteria": [
    "type:sparkeks"
  ],
  "cluster_criteria": [
    "data:prod"
  ],
  "context": {
    "query": "SELECT * from table limit 10;"
  }
}
```

## Result Format

If `return_result` is true, the plugin fetches results from S3 and returns them in tabular format:

```json
{
  "columns": [
    {"name": "col1", "type": "string"},
    {"name": "col2", "type": "int"}
  ],
  "data": [
    ["foo", 1],
    ["bar", 2]
  ]
}
```

## Testing

- **Local Docker**: Set AWS credentials in your environment or docker-compose.yml.
- **ECS Production**: The plugin uses the ECS task role for AWS authentication.

## Notes

- To use custom SparkApplication templates, provide the file path in `spark_application_file`.
- S3 URIs must be accessible by the Spark job and the plugin.
- For troubleshooting, check logs in the specified S3 logs URI.
