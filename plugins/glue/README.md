# 🍯 Glue Plugin

The **Glue Plugin** enables Heimdall to query the AWS Glue Data Catalog to retrieve metadata about a specific table. It helps you integrate Glue metadata queries directly into your orchestration workflows.

---

## 🧩 Plugin Overview

* **Plugin Name:** `glue`
* **Execution Mode:** Sync
* **Use Case:** Fetching AWS Glue table metadata for auditing, validation, or downstream processing

---

## ⚙️ Defining a Glue Command

```yaml
- name: glue-metadata-0.0.1
  status: active
  plugin: glue
  version: 0.0.1
  description: Query AWS Glue catalog for table metadata
  context:
    catalog_id: 123456789012
  tags:
    - type:glue-query
  cluster_tags:
    - type:localhost
```

* `catalog_id` is the AWS Glue catalog identifier (optional; defaults to AWS account ID if omitted).

---

## 🚀 Submitting a Glue Job

Specify the Glue table name to fetch metadata for in the job context:

```json
{
  "name": "fetch-glue-table-metadata",
  "version": "0.0.1",
  "command_criteria": ["type:glue-query"],
  "cluster_criteria": ["data:local"],
  "context": {
    "table_name": "my_database.my_table"
  }
}
```

---

## 📊 Returning Job Results

The plugin returns raw metadata as a JSON string containing the table schema and properties.

Retrieve results via:

```
GET /api/v1/job/<job_id>/result
```

---

## 🧠 Best Practices

* Use appropriate IAM permissions for Glue Catalog access.
* Validate the fully qualified `table_name` to avoid query errors.
* Use the metadata for auditing, documentation, or dynamic pipeline decisions.
