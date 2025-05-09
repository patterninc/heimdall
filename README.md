# Heimdall

**Heimdall** is a lightweight, pluggable data orchestration and job execution platform that abstracts complex data infrastructure from clients while offering a secure and consistent API for submitting and managing jobs.

Originally inspired by [Netflix Genie](https://github.com/Netflix/genie), Heimdall extends the architecture to support:

* 🔌 **Pluggable commands**
* ⚙️ **Job queuing**
* 📡 **Synchronous and asynchronous execution**

---

## ✨ Key Features

* 🔁 **Sync & Async Job Execution**
* 🧩 **Plugin-Based Execution Framework**: Shell, Glue, Snowflake, Spark, DynamoDB, and Ping
* 📬 **REST API** for programmatic access
* 🌍 **Web UI** for visual management
* 🔐 **Secure orchestration without credential leakage**
* 🧠 **Dynamic routing based on command / cluster criteria**
* 📦 **Configurable or self-registering clusters (future)**

---

## 🖥️ UI + API Access

Heimdall includes a web interface running alongside the API.

* **API**: `http://localhost:9090/api/v1`
* **Web UI**: `http://localhost:9090/ui`

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone git@github.com:patterninc/heimdall.git
cd heimdall
```

### 2. Start Heimdall

Ensure you have Docker or an compatible alternative installed.

```bash
docker compose up --build -d
```

This starts:

* The Heimdall server on port `9090`
* The database and other dependencies

### 3. Submit a Test Job (Example: Ping Plugin)

```bash
curl -X POST -H "X-Heimdall-User: test_user" -H "Content-Type: application/json" \
     -d '{
           "name": "ping-test",
           "version": "0.0.1",
           "context": {},
           "command_criteria": ["type:ping"],
           "cluster_criteria": ["type:localhost"]
         }' \
     http://127.0.0.1:9090/api/v1/job
```

### 4. Monitor the Job

Use the UI (`http://127.0.0.1:9090/`) or the following endpoints:

```bash
# Job status
GET /api/v1/job/<job_id>

# Job stdout
GET /api/v1/job/<job_id>/stdout

# Job stderr
GET /api/v1/job/<job_id>/stderr
```

---

## 🔌 Supported Plugins

Heimdall supports a growing set of pluggable command types:

| Plugin      | Description                            | Execution Mode |
| ----------- | -------------------------------------- | -------------- |
| `ping`      | Basic plugin used for testing          | Sync or Async  |
| `shell`     | Shell command execution                | Sync or Async  |
| `glue`      | Pulling Iceberg table metadata         | Sync or Async  |
| `dynamodb`  | DynamoDB read operation                | Sync or Async  |
| `snowflake` | Query execution in Snowflake           | Async          |
| `spark`     | SparkSQL query execution on EMR on EKS | Async          |

---

## 🧬 Core Concepts

### **Command**

Defines a reusable unit of work with associated tags and plugin logic.

### **Cluster**

An execution environment abstracted from its physical form. It can represent localhost, EMR, Kubernetes, a DB, a piece of your infrastructure that has context and a name, etc.

### **Job**

The orchestration request. It combines:

* Command criteria
* Cluster criteria
* Execution context

Heimdall dynamically selects the best command-cluster pair based on these criteria.

---

## ⚙️ Configuration

Initially, Commands and Clusters are configured via a static config file (see [config.yml](https://github.com/patterninc/heimdall/blob/main/configs/local.yaml)). Heimdall is evolving toward support for:

* Self-registering clusters
* Health-based routing
* API-based dynamic configuration

---

## 🔁 Command & Cluster Matching Logic

1. **Commands**: Must be active and match all tags in `command_criteria`.
2. **Compatible Clusters**: Found via the command’s own `cluster_criteria`.
3. **Final Selection**:

   * Filters clusters using the job’s `cluster_criteria`.
   * If multiple pairs match, one is selected randomly (a capability for custom "routing" is in works and will be represented as a plugin).
   * If no match, the job fails with a detailed error.

---

## 🔐 Security by Design

Heimdall removes the need for:

* Embedding credentials in user environments
* Direct user and services access to infrastructure

It centralizes execution logic, logging, and auditing—all accessible via API or UI.

---

## 📦 API Overview

| Endpoint                      | Description                    |
| ----------------------------- | ------------------------------ |
| `POST /api/v1/job`            | Submit a job                   |
| `GET /api/v1/job/<id>`        | Get job details                |
| `GET /api/v1/job/<id>/status` | Check job status               |
| `GET /api/v1/job/<id>/stdout` | Get stdout for a completed job |
| `GET /api/v1/job/<id>/stderr` | Get stderr for a completed job |
| `GET /api/v1/job/<id>/result` | Get job's result               |
| `GET /api/v1/commands`        | List configured commands       |
| `GET /api/v1/clusters`        | List configured clusters       |

---

## 👥 Credits

**Heimdall** was created at **Pattern, Inc** by Stan Babourine, with contributions from Will Graham, Gaurav Warale and Josh Diaz.
