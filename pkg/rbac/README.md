# RBAC Module

The RBAC (Role-Based Access Control) module provides authorization capabilities for Heimdall, specifically designed to integrate with Apache Ranger for fine-grained access control over SQL resources.

## Overview

This module enables Heimdall to:
- Authorize SQL queries based on Ranger policies
- Support table-level, schema-level, and catalog-level access control
- Handle user groups and permissions
- Automatically sync policies from Ranger

## Architecture

The module consists of several key components:

### Core Interfaces

- **`RBAC`**: Main interface for access control providers
### Apache Ranger Integration

The module currently supports Apache Ranger as the primary RBAC provider through:

- **`ApacheRanger`**: Main implementation of RBAC interface
- **`Policy`**: Represents Ranger policies with resources and permissions
- **`User`** and **`Group`**: Represent Ranger users and groups

## Configuration

### YAML Configuration Example

```yaml
rbacs:
  - type: apache_ranger
    name: my-ranger
    service_name: my_service
    sync_interval_in_minutes: 30
    client:
      url: https://ranger.example.com
      username: admin
      password: secret
    parser:
      type: trino
      default_catalog: hive
```
### Configuration Parameters

- **`type`**: RBAC provider type (`apache_ranger`)
- **`name`**: Unique identifier for this RBAC instance
- **`service_name`**: Ranger service name to fetch policies from
- **`sync_interval_in_minutes`**: How often to sync policies from Ranger
- **`client`**: Ranger connection configuration
- **`parser`**: SQL parser configuration for query analysis

### YAML Cluster Configuration

To enable RBAC in your cluster configuration, add the `rbacs` section and specify the RBAC provider names you want to use. This allows you to define a chain of permission providers.

Example:

```yaml
cluster:
    name: my-cluster
    rbacs:
        - my-ranger
        - another-rbac-provider
```

The `rbacs` list enables multiple RBAC providers to be evaluated in order, allowing flexible and layered access control.


## Usage

### Initialization

```go
import "github.com/patterninc/heimdall/pkg/rbac"

// Parse configuration
var rbacs rbac.RBACs
err := yaml.Unmarshal(configData, &rbacs)

// Initialize RBAC providers
ctx := context.Background()
for _, rbac := range rbacs {
    err := rbac.Init(ctx)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Access Control

```go
// Check if user has access to execute a query
user := "john.doe"
query := "SELECT * FROM catalog.schema.table"

hasAccess, err := rbac.HasAccess(user, query)
if err != nil {
    log.Error("Error checking access:", err)
    return
}

if !hasAccess {
    log.Info("Access denied for user:", user)
    return
}

// Execute query...
```

## Features

### Supported SQL Actions

The module supports fine-grained control over SQL operations:

- `SELECT`: Read data from tables
- `INSERT`: Insert data into tables
- `UPDATE`: Update existing data
- `DELETE`: Delete data from tables
- `CREATE`: Create new objects (tables, schemas, etc.)
- `DROP`: Drop existing objects
- `ALTER`: Modify existing objects
- `USE`: Use/switch to a schema or catalog
- `GRANT`/`REVOKE`: Manage permissions
- `SHOW`: Show system information
- `IMPERSONATE`: Act as another user
- `EXECUTE`: Execute procedures/functions

### Resource Matching

Policies support wildcard patterns for flexible resource matching:

- `*`: Matches any characters
- `?`: Matches single character
- Regular expressions for complex patterns

### Policy Types

- **Allow Policies**: Grant specific permissions to users/groups
- **Deny Policies**: Explicitly deny permissions (takes precedence)
- **Exceptions**: Override allow/deny policies for specific cases

### Automatic Synchronization

- Policies are automatically synced from Ranger at configured intervals
- Users and groups are kept up-to-date
- Background goroutine handles sync without blocking operations

## API Reference

### RBAC Interface

```go
type RBAC interface {
    Init(ctx context.Context) error
    HasAccess(user string, query string) (bool, error)
    GetName() string
}
```


## Error Handling

The module handles various error scenarios:

- **Network failures**: Retries and graceful degradation
- **Invalid policies**: Logs warnings and skips malformed policies
- **Unknown users**: Returns access denied for unknown users
- **Parsing errors**: Returns detailed error messages for invalid queries

## Performance Considerations

- **Caching**: User permissions are cached in memory for fast access
- **Batch requests**: API calls use pagination for large datasets
- **Background sync**: Policy updates don't block query processing
- **Regex compilation**: Resource patterns are pre-compiled for efficiency

## Security

- **Basic Authentication**: Secure communication with Ranger
- **Case-insensitive usernames**: Consistent user matching
- **Deny-by-default**: Unknown users/resources are denied access
- **Audit logging**: All access decisions are logged

## Troubleshooting

### Common Issues

1. **User not found in policies**
   - Ensure user exists in Ranger
   - Check group memberships
   - Verify sync is working

2. **Access denied unexpectedly**
   - Check policy resource patterns
   - Verify deny policies aren't blocking access
   - Review query parsing results

3. **Sync failures**
   - Verify Ranger connectivity
   - Check authentication credentials
   - Review Ranger service configuration

### Debugging

Enable verbose logging to see detailed access control decisions:

```go
log.SetLevel(log.DebugLevel)
```

This will show policy matches, resource patterns, and access decisions for each query.