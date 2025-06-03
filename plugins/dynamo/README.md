# 📦 DynamoDB Plugin

The **DynamoDB Plugin** enables Heimdall to run read-only PartiQL queries against AWS DynamoDB tables. It securely connects using AWS credentials or an optional assumed IAM role, supports pagination, and returns structured query results.

🔒 **Read-Only:** This plugin only supports data retrieval via PartiQL — no writes or modifications.

---

## 🧩 Plugin Overview

* **Plugin Name:** `dynamo`
* **Execution Mode:** Sync
* **Use Case:** Querying DynamoDB tables using PartiQL from within Heimdall workflows

---

## ⚙️ Defining a DynamoDB Command

### 📝 Command Definition

```yaml
- name: dynamo-0.0.1
  status: active
  plugin: dynamo
  version: 0.0.1
  description: Read data using PartiQL
  tags:
    - type:dynamodb
  cluster_tags:
    - type:dynamodb
```

### 🔐 Cluster Definition (AWS Authentication)

```yaml
- name: dynamo-0.0.1
  status: active
  version: 0.0.1
  description: AWS DynamoDB
  context:
    role_arn: arn:aws:iam::123456789012:role/HeimdallDynamoQueryRole
  tags:
    - type:dynamodb
    - data:prod
```

* `role_arn` is optional. If provided, Heimdall will assume this IAM role to authenticate requests.

---

## 🚀 Submitting a DynamoDB Job

Define the PartiQL query and optional result limit in the job context:

```json
{
  "name": "list-items-job",
  "version": "0.0.1",
  "command_criteria": ["type:dynamo"],
  "cluster_criteria": ["data:prod"],
  "context": {
    "query": "SELECT * FROM my_table WHERE category = 'books'",
    "limit": 100
  }
}
```

---

## 📊 Returning Job Results

The plugin paginates through query results automatically, returning a structured output like:

```json
{
  "columns": [
    {"name": "id", "type": "string"},
    {"name": "category", "type": "string"},
    {"name": "price", "type": "float"}
  ],
  "data": [
    ["123", "books", 12.99],
    ["124", "books", 15.50]
  ]
}
```

Retrieve results via:

```
GET /api/v1/job/<job_id>/result
```

---

## 🧠 Best Practices

* Use IAM roles with **least privilege** for security.
* Test queries in AWS console or CLI before running via Heimdall.
* Avoid large result sets by leveraging the `limit` parameter.
* Validate job input to prevent malformed PartiQL queries.
