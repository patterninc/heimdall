# ⚡ PostgreSQL Plugin

The **PostgreSQL Plugin** enables Heimdall to run SQL queries on configured PostgreSQL databases. It supports direct SQL, SQL files, batch execution, and both synchronous and asynchronous modes.

---

## 🧩 Plugin Overview

* **Plugin Name:** `postgres`
* **Execution Modes:** Synchronous (return_result: true) and Asynchronous (return_result: false)
* **Use Case:** Running SQL queries (single or batch) against PostgreSQL databases

---

## ⚙️ Defining a Postgres Command

A Postgres command can specify execution mode and other preferences. Example:

```yaml
  - name: postgres-0.0.1
    status: active
    plugin: postgres
    version: 0.0.1
    description: Execute queries against PostgreSQL databases
    tags:
      - type:postgres
    cluster_tags:
      - type:postgres
      - data:local
```

---

## 🖥️ Cluster Configuration

Each Postgres cluster must define a `connection_string`:

```yaml
  - name: postgres
    status: active
    version: 0.0.1
    description: PostgreSQL Production Database
    context:
      connection_string: "postgresql://user:password@host:port/database"
    tags:
      - type:postgres
      - data:local
```

---

## 🚀 Submitting a Postgres Job

A Postgres job provides the SQL query to be executed, and can specify execution mode:

```json
{
  "name": "run-pg-query",
  "version": "0.0.1",
  "command_criteria": [
    "type:postgres"
  ],
  "cluster_criteria": [
    "data:local"
  ],
  "context": {
    "query": "select * from employees limit 10;",
    "return_result": true
  }
}
```

---

## 📦 Job Context & Runtime

The Postgres plugin handles:

* Executing single or multiple SQL statements (batch)
* Supporting both direct query strings and SQL file execution
* Synchronous mode (`return_result: true`): returns query results, only one query allowed
* Asynchronous mode (`return_result: false`): executes all queries, returns success or error

### Job Context Example

```yaml
query: SELECT * FROM my_table # Required - SQL query to execute or path to .sql file
return_result: true # Optional - Whether to return query results (default: false)
```

### Cluster Context Example

```yaml
connection_string: postgresql://user:password@host:port/database # Required
```

---

## 📊 Returning Job Results

If enabled in the environment, Heimdall exposes query results via:

```
GET /api/v1/job/<job_id>/result
```

---

## 🧠 Best Practices

* Use synchronous mode for SELECT queries where results are needed
* Use asynchronous mode for DDL/DML or batch operations
* Always secure your connection strings and database credentials
* Use `type:postgres` tags to isolate command and cluster matching
* Use SQL files for large or complex batch operations
