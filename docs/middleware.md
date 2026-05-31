# Middleware

Middleware can be global (`server.middlewares`) or local (`routes[].middlewares`). Local middleware runs after routing and before checks.

Supported middleware:
- `cors`
- `recovery`
- `rate_limit`
- `logging`
- `request_id`
- `real_ip`
- `timeout`
- `metric`
- `inject`

## cors

```yaml
middlewares:
  - type: cors
    config:
      allowed:
        origin:
          - http://example.local
        method:
          - GET
          - POST
        header:
          - Content-Type
          - Authorization
```

`allowed.origin` is required. `method` and `header` can be empty, but the CORS response will not include those permissions.

## recovery

```yaml
middlewares:
  - type: recovery
```

Recovers from panics and returns HTTP 500.

## rate_limit

```yaml
middlewares:
  - type: rate_limit
    config:
      limit: 50
      window: 1m
```

Limits the number of requests within the configured `window`.

## logging

```yaml
middlewares:
  - type: logging
```

Logs basic request information.

## request_id

```yaml
middlewares:
  - type: request_id
```

Adds a request id to context/headers.

## real_ip

```yaml
middlewares:
  - type: real_ip
```

Detects the client IP from proxy headers when they are present.

## timeout

```yaml
middlewares:
  - type: timeout
    config:
      duration: 2s
```

Limits request processing time.

## metric

```yaml
middlewares:
  - type: metric
```

Collects simple request metrics.

## inject

```yaml
middlewares:
  - type: inject
    config:
      context:
        service_name: gateway
        version: v1
```

Adds literal values to `request.Context()`. This middleware currently does not perform an internal request and does not use Store.
