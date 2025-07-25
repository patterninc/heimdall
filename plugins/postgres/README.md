# PostgreSQL Plugin

This plugin provides an interface to PostgreSQL databases with support for direct SQL queries, SQL files, and batch execution.

## Features

- Execute single or multiple SQL statements
- Support for both direct query strings and SQL file execution
- Synchronous (return_result: true) and asynchronous (return_result: false) execution modes
- Error reporting with line number and query for batch execution
- Transaction support for batch execution

## Configuration

### Cluster Context

```yaml
connection_string: postgresql://user:password@host:port/database # Required
```

### Job Context

```yaml
query: SELECT * FROM my_table # Required - SQL query to execute or path to .sql file
return_result: true # Optional - Whether to return query results (default: false)
```

#### Execution Modes

1. **Return Results Mode** (`return_result: true`):
   - Only a single query is allowed (trailing semicolon is fine)
   - Returns query results as structured data
   - Fails if multiple queries are provided

2. **Execute Only Mode** (`return_result: false`):
   - Multiple queries separated by `;` are allowed
   - Executes all queries in order, stops and fails on first error
   - Returns error message with line number and query if any query fails
   - Returns a success message if all queries succeed

## Usage

```yaml
# Example 1: SELECT query with results
- name: postgresql-select
  description: Execute a SELECT query and return results
  command: postgres
  cluster: postgres-cluster
  context:
    query: SELECT * FROM my_table WHERE status = 'active'
    return_result: true

# Example 2: Execute an UPDATE statement
- name: postgresql-update
  description: Execute an UPDATE statement
  command: postgres
  cluster: postgres-cluster
  context:
    query: UPDATE my_table SET status = 'inactive' WHERE last_active < '2025-01-01'
    return_result: false


# Example 3: Execute multiple statements in a batch
- name: postgresql-batch
  description: Execute multiple SQL statements in a batch
  command: postgres
  cluster: postgres-cluster
  context:
    query: |
      INSERT INTO orders (id, customer_id, amount) VALUES ('ORD001', 'CUST123', 299.99);
      UPDATE inventory SET stock = stock - 1 WHERE product_id = 'PROD456';
    return_result: false
```
