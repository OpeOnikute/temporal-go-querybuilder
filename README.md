# Temporal Golang Query Builder

This is a simple query builder for Temporal Visibility API queries. It will live here until some consensus is reached in https://github.com/temporalio/features/issues/568.

## Usage
The Visibility API queries are simple strings that look like this:
```
ExecutionStatus='Running' AND WorkflowType='TestMe' 
```

Instead of concatenating strings (and other types) to make this queries, you can achieve the same using this package:
```
q := temporal.QueryBuilder{}
// assume that `temporal` is an internal package with re-usable constants
q.StartQuery(temporal.SearchAttrExecutionStatus, '=', temporal.ExecStatusRunning)
q.And(temporal.SearchAttrWorkflowType, '=', workflowName)
q.Encode()
```

This design is inspired by [net/url#values](https://pkg.go.dev/net/url#Values).

## Supported Logical Operators
- Or
- And
- Between
- In
- StartsWith