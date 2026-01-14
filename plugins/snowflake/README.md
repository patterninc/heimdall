# ❄️ Snowflake Plugin

The **Snowflake Plugin** enables Heimdall to execute **SQL queries** against a configured Snowflake instance using **key-pair authentication**. This allows seamless integration of Snowflake into your orchestration flows.

⚠️ **Note:** The plugin supports **a single SQL statement per job**.
Multi-statement execution is ***currently not supported***, but is planned for future updates.

---

## 🧩 Plugin Overview

* **Plugin Name:** `snowflake`
* **Execution Mode:** Async
* **Use Case:** Executing data queries or transformations on Snowflake

---

## ⚙️ Defining a Snowflake Command

You must define both a **command** and a **cluster**. The command represents the logical query job, while the cluster provides Snowflake connection details.

### 🔹 Command Configuration

```yaml
- name: snowflake-0.0.1
  status: active
  plugin: snowflake
  version: 0.0.1
  description: Query user metrics from Snowflake
  context:
    role: DATA_ENGINEER
  tags:
    - type:snowflake
  cluster_tags:
    - type:snowflake
```

### 🔸 Cluster Configuration

The cluster configuration represents **identity** (who/where), not **permissions** (what they can do).
Role is intentionally excluded to prevent configuration duplication across clusters.

```yaml
- name: snowflake-prod-cluster
  status: active
  version: 0.0.1
  description: Production Snowflake cluster
  context:
    account: myorg-account-id
    user: my-snowflake-user
    database: MY_DB
    warehouse: MY_WAREHOUSE
    private_key: /etc/keys/snowflake-private-key.p8
  tags:
    - type:snowflake
    - data:prod
```

> The `private_key` field must point to a valid **PKCS#8** PEM-formatted file accessible from the execution environment.

---

## 🚀 Submitting a Snowflake Job

Jobs must include a single SQL statement via the `context.query` field.

```json
{
  "name": "country-count-report",
  "version": "0.0.1",
  "command_criteria": ["type:user-metrics"],
  "cluster_criteria": ["type:snowflake-prod"],
  "context": {
    "query": "SELECT country, COUNT(*) AS users FROM user_data GROUP BY country"
  }
}
```

🔹 The plugin will execute this query using the Snowflake configuration defined in the matched cluster.

> Avoid semicolon-separated or batch queries—only **a single statement per job** is supported at this time.

---

## 🎭 Role Configuration

The Snowflake role determines what permissions the connection has. Role is **decoupled from cluster configuration** and defined at the **command level** for security:
- Prevent configuration duplication across clusters
- **Prevent per-job role escalation** - users cannot override roles in job submissions

### ⚠️ Default Behavior

If no role is specified in the command context, Snowflake uses the **user's default role**.

---

## 📊 Returning Query Results

If your query returns data, Heimdall captures it as structured output accessible via:

```
GET /api/v1/job/<job_id>/result
```

✅ Example format:

```json
{
  "columns": [
    {"name": "country", "type": "string"},
    {"name": "users", "type": "int"}
  ],
  "data": [
    ["USA", 1500],
    ["Canada", 340]
  ]
}
```

---

## 🔐 Authentication & Security

* Auth is performed using Snowflake's **key-pair authentication**.
* Ensure private keys are stored securely and only readable by the runtime environment.

---

## 🧠 Best Practices

* Write simple, single-statement queries for now.
* Keep secrets (like private keys) out of job context—store them securely in the cluster configuration.
* Use cluster and command tags to enforce environment separation (e.g., `prod`, `dev`, etc.).
* Plan for upcoming multi-statement support by structuring your queries modularly where possible.
