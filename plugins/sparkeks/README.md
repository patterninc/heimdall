# Spark EKS Plugin

The Spark EKS plugin enables submitting Spark jobs to an AWS EKS cluster using the Kubeflow Spark Operator. It supports IAM role assumption, S3 integration for queries/results/logs, and custom Spark application templates.

## Features

- Submit Spark jobs to EKS via Kubeflow Spark Operator
- Upload SQL queries to S3 and fetch results from S3
- Collect Spark pod logs and store them in S3
- Supports custom SparkApplication YAML templates
- IAM role assumption for cross-account EKS access
- Configurable Spark job resources and properties
- **JAR / `--class` entrypoints** — run a JVM (Scala/Java) application by main class, alongside the default Python SQL wrapper

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

## Entrypoints: SQL wrapper vs JAR

The plugin picks an entrypoint from the **file extension of the command's `wrapper_uri`**. The wrapper
file itself is always submitted as `MainApplicationFile`; only `Type`, `MainClass`, and `Arguments`
differ between the two.

| `wrapper_uri` ends in | Entrypoint | `Spec.Type` | `Spec.MainClass` |
|---|---|---|---|
| `.py` | Python SQL wrapper (default, pre-existing behavior) | `Python` | not set |
| `.jar` | JVM JAR, run by main class | `Scala` or `Java` | `entry_point` |
| anything else | falls back to the SQL wrapper | `Python` | not set |

### JAR configuration

Set on the **job context** under `parameters`:

| Field | Required | Description |
|---|---|---|
| `entry_point` | yes (for JAR) | Fully-qualified main class, e.g. `com.pattern.chipmunk.Writer`. Becomes `--class`. |
| `application_type` | no | `Scala` or `Java`, case-insensitive. Defaults to `Scala` when empty or unrecognized. |

```json
{
  "query": "SELECT key, value FROM my_table",
  "parameters": {
    "entry_point": "com.pattern.chipmunk.Writer",
    "application_type": "Scala",
    "properties": {
      "spark.chipmunk.output.path": "s3://pattern-dl/chipmunk/collections/my_collection/v1"
    }
  },
  "return_result": false
}
```

with the command context pointing `wrapper_uri` at the JAR:

```json
{
  "wrapper_uri": "s3://mybucket/chipmunk-assembly.jar"
}
```

### Argument contract

Both entrypoints receive the **same positional arguments**, so a JAR's `main()` must accept the
identical shape the SQL wrapper does:

```
[appName, queryURI, user]              # return_result = false
[appName, queryURI, user, resultURI]   # return_result = true
```

`queryURI` is an `s3a://` URI to the uploaded query file — the application reads the SQL **from that
URI**, it is not passed inline. There is no positional slot for job-specific inputs: the job context's
`arguments` field is ignored by both strategies. Pass such inputs as SparkConf `properties` instead
(chipmunk uses `spark.chipmunk.output.path`).

### `s3://` vs `s3a://` in properties

Values in `properties` (cluster and job) are passed through **verbatim**, with one exception: keys read
through Hadoop's `FileSystem` abstraction are rewritten `s3://` → `s3a://`. Today that is:

- `spark.kubernetes.driver.podTemplateFile`
- `spark.kubernetes.executor.podTemplateFile`

Everything else keeps the literal scheme you set, because application code typically reads these via
the AWS SDK, which does not understand `s3a://`. The plugin-managed URIs (`queryURI`, `resultURI`,
`event_log_uri`, `wrapper_uri`) are always rewritten to `s3a://`.

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
