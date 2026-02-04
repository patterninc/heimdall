# DBT Plugin for Heimdall

This plugin enables Heimdall to execute dbt commands against Snowflake.

## Overview

The dbt plugin allows Airflow DAGs to trigger dbt transformations through Heimdall, maintaining centralized credential management and job tracking.

## Features

- Execute dbt commands (`build`, `run`, `test`, `snapshot`)
- Support for selectors and model specifications
- Dynamic profiles.yml generation from cluster context
- Snowflake authentication via JWT (private key)
- Structured logging and error reporting

## Configuration

### Command Configuration

```yaml
commands:
  - name: dbt-0.0.1
    status: active
    plugin: dbt
    version: 0.0.1
    description: Run dbt commands against Snowflake
    context:
      dbt_project_path: /opt/dbt/analytics_dbt_poc
      dbt_profiles_dir: /opt/dbt/profiles
    tags:
      - type:dbt
    cluster_tags:
      - type:snowflake
```

### Cluster Configuration

Uses existing Snowflake cluster configuration:

```yaml
clusters:
  - name: snowflake-heimdall-wh
    context:
      account: pattern
      user: heimdall__app_user
      database: iceberg_db
      warehouse: heimdall_wh
      private_key: /etc/pattern.d/snowflake_key.pem
    tags:
      - type:snowflake
      - data:prod
```

## Job Context

When submitting a job from Airflow, provide the following context:

```json
{
  "command": "build",
  "selector": "hourly_models",
  "threads": 4,
  "target": "prod",
  "full_refresh": false
}
```

### Supported Parameters

- `command` (string): dbt command to run (`build`, `run`, `test`, `snapshot`)
- `selector` (string, optional): Selector name from selectors.yml
- `models` ([]string, optional): List of specific models to run
- `exclude` ([]string, optional): Models to exclude
- `threads` (int): Number of threads for parallel execution
- `full_refresh` (bool): Whether to perform a full refresh
- `vars` (map[string]string, optional): dbt variables
- `target` (string): Target environment (prod, dev)

## Authentication

The plugin uses Snowflake JWT authentication:

1. Reads private key from cluster context path
2. Generates profiles.yml with JWT authenticator
3. Executes dbt with proper credentials

## Example Usage from Airflow

```python
from pattern.operators.dbt import DbtOperator

dbt_task = DbtOperator(
    task_id="dbt_build_hourly",
    command="build",
    selector="hourly_models",
    threads=4,
    dag=dag
)
```

## Requirements

- dbt-core >= 1.7.0
- dbt-snowflake >= 1.7.0
- Valid Snowflake credentials in cluster configuration
