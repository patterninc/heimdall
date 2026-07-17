# 🗃️ StarRocks Plugin

The **StarRocks Plugin** enables Heimdall to execute SQL queries on configured StarRocks clusters using the [Arrow Flight SQL](https://arrow.apache.org/docs/format/FlightSql.html) protocol. It connects directly to a StarRocks FE's Arrow Flight SQL service and executes queries, returning results as native Arrow record batches.

---

## 🧩 Plugin Overview

* **Plugin Name:** `starrocks`
* **Use Case:** Running SQL queries against StarRocks clusters via Arrow Flight SQL, with optional result retrieval
* **Transport:** gRPC / Arrow Flight SQL (not the MySQL wire protocol)

---

## ⚙️ Defining a StarRocks Command

A StarRocks command defines authentication credentials for connecting to StarRocks clusters. The credentials are shared across all jobs using this command.

```yaml
  - name: starrocks-analytics
    status: active
    plugin: starrocks
    version: 0.0.1
    description: Execute StarRocks queries
    context:
      username: analytics_user
      password: secure_password
    tags:
      - type:starrocks
    cluster_tags:
      - type:starrocks
```

🔸 The command stores authentication credentials (`username` and `password`) used to perform the Arrow Flight SQL basic-auth handshake. StarRocks exchanges these for a bearer token which is attached to every subsequent request on the connection.

---

## 🖥️ Cluster Configuration

Each StarRocks cluster must define the Arrow Flight SQL `endpoint` (FE host and `arrow_flight_port`, default `9408`).

```yaml
  - name: starrocks-prod
    status: active
    version: 0.0.1
    description: Production StarRocks cluster
    context:
      endpoint: "starrocks-fe.company.com:9408"
      database: analytics
      use_tls: false
    tags:
      - type:starrocks
      - env:production
```

🔹 `endpoint` is the StarRocks FE's Arrow Flight SQL address (`host:port`). `database` is currently informational (set via the query itself, e.g. `use analytics; ...` or fully-qualified table names). `use_tls` controls whether the gRPC connection is established over TLS; StarRocks' Arrow Flight SQL service supports plaintext gRPC by default.

---

## 🚀 Submitting a StarRocks Job

A StarRocks job provides the SQL query to execute and result handling preferences.

```json
{
  "name": "user-analytics-query",
  "version": "1.0.0",
  "command_criteria": ["type:starrocks"],
  "cluster_criteria": ["env:production"],
  "context": {
    "query": "SELECT user_id, COUNT(*) AS events FROM user_events GROUP BY user_id",
    "return_result": true
  }
}
```

🔹 When `return_result` is `true`, the plugin runs the query via Flight SQL `GetFlightInfo`/`DoGet`, streaming Arrow record batches back and converting them into Heimdall's structured result format. When `false`, the statement is executed via Flight SQL's `ExecuteUpdate` (appropriate for `INSERT`/`UPDATE`/`DELETE`/DDL), and the result is a single `message` row reporting the number of affected rows.

---

## 📦 Job Context & Runtime

The StarRocks plugin handles:

* **Connection Management**: Dials the FE's Arrow Flight SQL gRPC endpoint and performs the Basic → Bearer token handshake
* **Query Execution**: Executes SQL via Arrow Flight SQL `Execute`/`ExecuteUpdate`
* **Type Handling**: Converts Arrow column types (ints, floats, decimals, strings, booleans, dates/timestamps) into Heimdall's result column types
* **Result Collection**: `GetFlightInfo` returns one endpoint per BE holding a slice of the result. Rather than proxying every row back through the FE, the plugin dials each distinct BE location directly (reusing the bearer token from the FE handshake — no re-authentication needed) and redeems tickets with `DoGet` in parallel across goroutines. Endpoints that advertise `LocationReuseConnection` or the FE's own address are still served over the existing FE connection. Direct BE connections are cached per query and closed once all endpoints are drained.

### Supported Arrow → Result Type Mapping

| Arrow Type | Result Type |
|---|---|
| `INT8`, `INT16`, `INT32`, `UINT8`, `UINT16` | `int` |
| `INT64`, `UINT32`, `UINT64` | `long` |
| `FLOAT32` | `float` |
| `FLOAT64`, `DECIMAL128`, `DECIMAL256` | `double` |
| `STRING`, `LARGE_STRING`, `BINARY`, `LARGE_BINARY` | `string` |
| `BOOL` | `boolean` |
| `DATE32`, `DATE64`, `TIMESTAMP`, `TIME32`, `TIME64` | `string` (formatted) |
| `LIST`, `LARGE_LIST`, `MAP`, `STRUCT` | `string` (JSON-encoded) |

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
    {"name": "user_id", "type": "long"},
    {"name": "events", "type": "long"}
  ],
  "data": [
    [12345, 156],
    [67890, 203],
    [11111, 89]
  ]
}
```
