# 🔥 Spark Plugin

The **Spark Plugin** enables Heimdall to submit SparkSQL batch jobs on AWS EMR on EKS clusters. It uploads SQL queries to S3, runs them via EMR Containers, and optionally returns query results in Avro format.

---

## 🧩 Plugin Overview

* **Plugin Name:** `spark`
* **Execution Mode:** Asynchronous batch jobs
* **Use Case:** Running SparkSQL queries with configurable properties on EMR on EKS clusters

---

## ⚙️ Defining a Spark Command

A Spark command requires a SQL `query` in its job context and optionally job-specific Spark properties. The plugin uploads queries and reads results from S3 paths configured in the command context.

```yaml
  - name: spark-sql-3.5.3
    status: active
    plugin: spark
    version: 3.5.3
    description: Run a SparkSQL query
    context:
      queries_uri: s3://bucket/spark/queries
      results_uri: s3://bucket/spark/results
      logs_uri: s3://bucket/spark/logs
      <!--
        - For SQL Python usage, provide the path to the `.py` file (e.g., s3://bucket/contrib/spark/spark-sql-s3-wrapper.py).
        - For Spark Applications(JAR) usage, set this to the path of the JAR file.
      -->
      wrapper_uri: s3://bucket/contrib/spark/spark-sql-s3-wrapper.py
      properties:
        spark.executor.instances: "1"
        spark.executor.memory: "500M"
        spark.executor.cores: "1"
        spark.driver.cores: "1"
    tags:
      - type:sparksql
    cluster_tags:
      - type:spark
```

🔸 This defines the S3 URIs for queries, results, logs, and [python wrapper script](https://github.com/patterninc/heimdall/blob/main/configs/spark-sql-s3-wrapper.py) used for queries that return results. Additional Spark properties can be specified globally or per job. Job properties take precedence over properties identified in the command.

---

## 🖥️ Cluster Configuration

Clusters must specify their EMR execution role ARN, EMR release label, and optionally an IAM RoleARN to assume:

```yaml
  - name: emr-on-eks
    status: active
    version: 3.5.3
    description: EMR on EKS cluster (7.6.0)
    context:
    execution_role_arn: arn:aws:iam::123456789012:role/EMRExecutionRole
    emr_release_label: emr-6.5.0-latest
    role_arn: arn:aws:iam::123456789012:role/AssumeRoleForEMR
      properties:
        spark.driver.memory: "2G"
        spark.sql.catalog.glue_catalog: "org.apache.iceberg.spark.SparkCatalog"
    tags:
      - type:spark
      - data:prod
```

---

## 🚀 Submitting a Spark Job

A typical Spark job includes a SQL query and optional Spark properties:

```json
{
  "name": "run-my-query",
  "version": "0.0.1",
  "command_criteria": ["type:sparksql"],
  "cluster_criteria": ["data_prod"],
  "context": {
    // For SQL jobs, specify the "query" field. For Spark jobs that execute a custom JAR, use "arguments" and parameters."entry_point".
    "query": "SELECT * FROM my_table WHERE dt='2023-01-01'",
    "arguments": ["SELECT 1", "s3:///"],
    "parameters": {
      // All values in "properties" are passed as `--conf` Spark submit parameters. Defaults from the command or cluster properties are merged in.
      "properties": {
        "spark.executor.memory": "4g",
        "spark.executor.cores": "2"
      },
      // The value of this property will be passed as the `--class` argument in Spark submit parameters,
      // specifying the main class to execute in your Spark application.
      "entry_point": "com.your.company.ClassName"
    },
    "return_result": true
  }
}
```

🔹 The job uploads the query to S3, submits it to EMR on EKS, and polls until completion. If `return_result` is true, results are fetched from S3 in Avro format.

---

## 📦 Job Context & Runtime

The Spark plugin handles:

* Uploading the SQL query file to the configured `queries_uri` in S3
* Submitting the job with the specified properties
* Polling the EMR job status until completion or failure
* Fetching results from the `results_uri` if requested
* Streaming job logs to configured `logs_uri`

---

## 📊 Returning Job Results

If `return_result` is set, Heimdall reads Avro-formatted query results from the results S3 path and exposes them via:

```
GET /api/v1/job/<job_id>/result
```

---

## 🧠 Best Practices

* Provide IAM roles with minimal necessary permissions for EMR and S3 access.
* Tune Spark properties carefully to optimize resource usage and performance.
* Use the `wrapper_uri` to customize job execution logic as you need it
* Monitor job states and handle failures gracefully in your workflows.
