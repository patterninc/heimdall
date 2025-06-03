# 🏓 Ping Plugin

The **Ping Plugin** is a sample command used for testing Heimdall’s orchestration flow. Instead of sending actual ICMP packets, it responds instantly with a predefined message — perfect for dry runs, plugin testing, or just checking your Heimdall wiring. 🚧

⚠️ **Testing Only:** This plugin is a *no-op*. It does **not** reach out to real hosts. Use it to verify that jobs run through Heimdall correctly.

---

## 🧩 Plugin Overview

* **Plugin Name:** `ping`
* **Execution Mode:** Sync
* **Use Case:** Testing job submission, validation, or plugin behavior without side effects

---

## ⚙️ Defining a Ping Command

You don’t need to specify much — just use the `ping` plugin and give it a name.

```yaml
- name: ping-0.0.1
  status: active
  plugin: ping
  version: 0.0.1
  description: Check Heimdall wiring
  tags:
    - type:ping
  cluster_tags:
    - type:localhost
```

🔹 When this job runs, Heimdall will simulate a ping and respond with a message like:

```
Hello, <calling user>!
```

---

## 🖥️ Cluster Configuration

Use a simple localhost cluster (or any compatible test target) to execute ping jobs:

```yaml
- name: localhost-0.0.1
  status: active
  version: 0.0.1
  description: Localhost
  tags:
    - type:localhost
    - data:local
```

---

## 🚀 Submitting a Ping Job

Here’s how to submit an example ping command via the Heimdall API:

```json
{
  "name": "ping-check-job",
  "version": "0.0.1",
  "command_criteria": ["type:ping"],
  "cluster_criteria": ["data:local"],
  "context": {}
}
```

🟢 This will run the ping plugin and instantly return result.

---

## 📊 Returning Job Results

The plugin returns this result:

```json
{
  "columns": [
    {"name": "message", "type": "string"}
  ],
  "data": [
    ["Hello, alice!"]
  ]
}
```

You can retrieve the result from:

```
GET /api/v1/job/<job_id>/result
```

---

## 🧠 Best Practices

* Use this plugin to **test your pipelines** before running real jobs.
* It’s great for **CI/CD checks**, plugin regression tests, or mocking command behavior.
* Don't forget: **no real pinging happens** — it's just a friendly "Hello!" 🎯
