# ⚡ Trino Plugin

The **Trino Plugin** enables Heimdall to run SQL queries on a configured Trino cluster. It sends SQL to a specified Trino endpoint and polls for results asynchronously, eliminating the need to implement pagination on the client.

---

## 🧩 Plugin Overview

* **Plugin Name:** `trino`
* **Execution Mode:** Asynchronous query execution (recommended)
* **Use Case:** Running Trino SQL queries against configured catalogs and endpoints

---

## ⚙️ Defining a Trino Command

A Trino command defines optional polling behavior. It does not require static SQL or properties, but can set polling preferences shared across jobs.

```yaml
  - name: trino-414
    status: active
    plugin: trino
    version: 414
    description: Run Trino queries
    context:
      poll_interval: 2000  # milliseconds between poll attempts
    tags:
      - type:trino
    cluster_tags:
      - type:trino
```

🔸 This command sets `poll_interval` to 2000 milliseconds (2 seconds) between polling attempts. Command-level settings apply to all jobs unless overridden by the job context.

---

## 🖥️ Cluster Configuration

Each Trino cluster must define a coordinator `endpoint`, and can optionally set a default `catalog`.

```yaml
  - name: trino-cluster-west
    status: active
    version: 414
    description: Trino cluster in us-west-2
    context:
      endpoint: https://trino.company.com:8443
      catalog: hive
    tags:
      - type:trino
      - region:us-west
```

🔹 The `endpoint` points to the Trino coordinator. If a `catalog` is defined here, it becomes the default for all jobs targeting this cluster unless overridden in future plugin versions.

---

## 🚀 Submitting a Trino Job

A Trino job provides the SQL query to be executed, and inherits command and cluster context.

```json
{
  "name": "run-trino-query",
  "version": "0.1.0",
  "command_criteria": ["type:trino"],
  "cluster_criteria": ["region:us-west"],
  "context": {
    "query": "SELECT order_id, total FROM orders WHERE order_date >= DATE '2024-01-01'"
  }
}
```

🔹 The job sends the SQL query to the Trino endpoint, then polls until execution is complete. Results are streamed back once ready—no client-side pagination needed.

---

## 📦 Job Context & Runtime

The Trino plugin handles:

* Sending the SQL query to the configured Trino cluster endpoint
* Polling for completion using the specified `poll_interval` (in milliseconds)
* Returning the full result set when the query completes
* Avoiding pagination by retrieving results in one asynchronous operation

---

## 📊 Returning Job Results

If enabled in the environment, Heimdall exposes query results via:

```
GET /api/v1/job/<job_id>/result
```

---

## 🧠 Best Practices

* Use the asynchronous mode to avoid blocking behavior and simplify job orchestration
* Configure Trino endpoints securely (HTTPS, access controls, etc.)
* Adjust `poll_interval` to reduce load on the Trino coordinator while keeping jobs responsive
* Use `type:trino` tags to isolate command and cluster matching
* Set `catalog` defaults in clusters to avoid redundancy in job configurations
