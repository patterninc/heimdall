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
| `entry_point` | yes (for JAR) | Fully-qualified main class, e.g. `com.org.customapplication.main`. Becomes `--class`. |
| `application_type` | no | `Scala` or `Java`, case-insensitive. Defaults to `Scala` when empty or unrecognized. |

```json
{
  "query": "SELECT key, value FROM my_table",
  "wrapper_uri": "s3://mybucket/application-assembly.jar",
  "parameters": {
    "entry_point": "com.org.customapplication.main",
    "application_type": "Scala",
    "properties": {
      //sparkeks properties
    }
  },  
  "return_result": false
}
```

### Arguments

`arguments` **extends** the managed argument list, it does not replace it.

| index | value |
|---|---|
| 0 | `app_name` |
| 1 | `query_uri` — `s3a://` path to the uploaded `query.sql`; the app reads the SQL from it |
| 2 | `user` — the authenticated caller; set `kyuubi.session.user` from this |
| 3 | `result_uri`, or `""` when `return_result` is false |
| 4+ | the job context's `arguments`, in order |

The indices are the same in both supportedruntimes: Scala `args(i)` and Python `sys.argv[i+1]` both map to
slot `i`.

#### Example 1 — no extras

```json
{
  "context": {
    "query": "SELECT * FROM my_table LIMIT 10;",
    "return_result": true
  }
}
```

The app receives `[app_name, query_uri, user, result_uri]`.

#### Example 2 — app-specific extras

An app needing its own output prefix passes it as an extra; the managed slots are untouched:

```json
{
  "context": {
    "query": "SELECT key, value FROM my_table",
    "arguments": [
      "s3://mybucket/output/v1"
    ],
    "parameters": {
      "entry_point": "com.org.customapplication.main",
      "application_type": "Scala"
    }
  }
}
```

The app receives `[app_name, query_uri, user, result_uri, "s3://mybucket/output/v1"]`

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
