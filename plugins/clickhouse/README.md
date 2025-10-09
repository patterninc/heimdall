# 🗃️ ClickHouse Plugin

The **ClickHouse Plugin** enables Heimdall to execute SQL queries on configured ClickHouse clusters. It connects directly to ClickHouse instances and executes queries, with support for parameterized queries and optional result collection.

---

## 🧩 Plugin Overview

* **Plugin Name:** `clickhouse`
* **Use Case:** Running SQL queries against ClickHouse databases with optional result retrieval

---

## ⚙️ Defining a ClickHouse Command

A ClickHouse command defines authentication credentials for connecting to ClickHouse clusters. The credentials are shared across all jobs using this command.

```yaml
  - name: clickhouse-analytics
    status: active
    plugin: clickhouse
    version: 24.8
    description: Execute ClickHouse queries
    context:
      username: analytics_user
      password: secure_password
    tags:
      - type:clickhouse
    cluster_tags:
      - type:clickhouse
```

🔸 The command stores authentication credentials (`username` and `password`) that will be used for all ClickHouse connections. These credentials are applied to all jobs targeting ClickHouse clusters.

---

## 🖥️ Cluster Configuration

Each ClickHouse cluster must define connection `endpoints` and optionally specify a default `database`.

```yaml
  - name: clickhouse-prod
    status: active
    version: 24.8
    description: Production ClickHouse cluster
    context:
      endpoints:
        - "clickhouse-node1.company.com:9000"
        - "clickhouse-node2.company.com:9000"
        - "clickhouse-node3.company.com:9000"
      database: analytics
    tags:
      - type:clickhouse
      - env:production
```

🔹 The `endpoints` array defines ClickHouse server addresses with ports. The optional `database` parameter sets the default database for connections to this cluster.

---

## 🚀 Submitting a ClickHouse Job

A ClickHouse job provides the SQL query to execute, optional parameters, and result handling preferences.

```json
{
  "name": "user-analytics-query",
  "version": "1.0.0",
  "command_criteria": ["type:clickhouse"],
  "cluster_criteria": ["env:production"],
  "context": {
    "query": "SELECT user_id, COUNT(*) AS events FROM user_events WHERE date >= {date:Date} AND user_type = {user_type:String} GROUP BY user_id",
    "params": {
        "date": "2024-01-01",
        "user_type": "premium"
    },
    "return_result": true
  }
}
```

🔹 The job executes the SQL query with the provided parameters and returns results if `return_result` is enabled. Parameters are safely bound to prevent SQL injection.

---

## 📦 Job Context & Runtime

The ClickHouse plugin handles:

* **Connection Management**: Establishes secure connections to ClickHouse clusters using provided credentials
* **Query Execution**: Executes SQL queries with parameter binding for security
* **Type Handling**: Properly handles ClickHouse data types including nullable and low cardinality variants
* **Result Collection**: Optionally collects and formats query results for API responses

### Supported ClickHouse Data Types

The plugin supports comprehensive ClickHouse type mapping:

| ClickHouse Type | Go Type | Nullable Support |
|----------------|---------|------------------|
| `UInt8`, `UInt16`, `UInt32`, `Int8`, `Int16`, `Int32`,   | `int`, | ✅ |
|  `UInt64`   `Int64`                        | `int64` | ✅ |
| `Float32`, `Float64`                       | `float32`, `float64` | ✅ |
| `String`, `FixedString`                    | `string` | ✅ |
| `Date`, `Date32`, `DateTime`, | `time.Time` | ✅ |
| `Decimal(P,S)`,     | `string` | ✅ |

🔸 The plugin automatically handles:
- **Nullable types**: `Nullable(String)` → `*string`
- **Low cardinality**: `LowCardinality(String)` → `string`
- **Complex wrappers**: `Nullable(LowCardinality(String))` → `*string`
- **Decimal variants**: `Decimal64(18,4)` → `string`

---

## 📊 Returning Job Results

When `return_result` is enabled, query results are available via:

```
GET /api/v1/job/<job_id>/result
```

Results are returned in structured format:

```json
{
  "columns": [
    {"name": "user_id", "type": "UInt64"},
    {"name": "events", "type": "UInt64"}
  ],
  "data": [
    [12345, 156],
    [67890, 203],
    [11111, 89]
  ]
}
```