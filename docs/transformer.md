# Transformer

A transformer is a route-level layer that modifies the incoming request before proxying. It can add or replace headers, query parameters, and JSON body fields.

It is the last route-specific step before the proxy layer.

## Supported Types

### header

Transforms request headers.

```yaml
transforms:
  header:
    X-User-ID: "{route:id}"
```

Existing header values are overwritten.

### body_fields

Transforms the request body.

```yaml
transforms:
  body_fields:
    user:
      id: "{route:id}"
```

Nested objects are supported. If the original body does not contain the nested object, it will be created.

Existing field values are overwritten.

Arrays are not supported yet.

### query_params

Adds or replaces query parameters.

```yaml
transforms:
  query_params:
    locale: "{query:locale}"
```

Existing query parameter values are overwritten.
