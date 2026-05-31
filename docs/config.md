# Configuration

The `config.yaml` file contains `server` and `routes`.

```yaml
server:
  middlewares: []
routes: []
```

## Server

```yaml
server:
  middlewares:
    - type: recovery
    - type: logging
```

`server.middlewares` defines global middleware applied to every route.

## Route

```yaml
routes:
  - path: /users/{id}
    method: GET
    middlewares: []
    checks: []
    transforms: {}
    upstream: {}
```

Fields:
- `path` - incoming route path. It can contain chi parameters, for example `/users/{id}`.
- `method` - incoming HTTP method.
- `middlewares` - middleware applied only to this route.
- `checks` - validations executed before transforms and proxying.
- `transforms` - request modifications applied before proxying.
- `upstream` - target service configuration.

## Upstream

```yaml
upstream:
  scheme: http
  host: api.local
  path: /users/{route:id}
  method: GET
```

Fields:
- `scheme` - `http`, `https`, or `ws`.
- `host` - host without scheme.
- `path` - upstream path. It can contain request templates.
- `method` - upstream method. For route upstreams, an empty value falls back to `route.method`.

## Resolver

Template format: `{source:key}`.

Available sources for the request renderer:
- `{context:key}`
- `{route:key}`
- `{query:key}`
- `{header:key}`

Available sources for the response renderer in `store`:
- `{header:key}`
- `{body:key}`

`{body:key}` reads only a top-level field from a JSON object in the response body. Nested paths and arrays are not supported yet.

## Checks

Supported checks match the registrations in `internal/builder/bootstrap.go`: `policy`, `header_required`, `ip_whitelist`, and `query_required`.

### policy

```yaml
checks:
  - type: policy
    config:
      transform:
        header:
          X-Request-ID: "{header:X-Request-ID}"
        query_params:
          locale: "{query:locale}"
      upstream:
        scheme: http
        host: policy.local
        path: /auth/{route:id}
        method: POST
      expected_status: 200
      store:
        token: "{header:X-Token}"
        user_id: "{body:user_id}"
```

`policy` performs an internal request to `upstream`, applies `transform` to that request, compares the response status with `expected_status`, and then stores response values through `store`.

### header_required

```yaml
checks:
  - type: header_required
    config:
      headers:
        - X-Request-ID
```

### ip_whitelist

```yaml
checks:
  - type: ip_whitelist
    config:
      ips:
        - 127.0.0.1
```

### query_required

```yaml
checks:
  - type: query_required
    config:
      query_params:
        - locale
```

## Store

```yaml
store:
  token: "{header:X-Token}"
  user_id: "{body:user_id}"
```

Store works with `http.Response`, not with the incoming `http.Request`. This makes it suitable only for checks or middleware that perform an internal request and receive a response.

## Transforms

```yaml
transforms:
  header:
    X-User-ID: "{route:id}"
  query_params:
    locale: "{query:locale}"
  body_fields:
    user:
      id: "{route:id}"
```

Types:
- `header` - sets request headers.
- `query_params` - sets query parameters.
- `body_fields` - merges fields into a JSON body. Objects are supported recursively, but arrays are not supported yet.

## Middleware

```yaml
middlewares:
  - type: cors
    config:
      allowed:
        origin:
          - http://example.local
        method:
          - GET
        header:
          - Authorization
```

Supported middleware: `cors`, `recovery`, `rate_limit`, `logging`, `request_id`, `real_ip`, `timeout`, `metric`, and `inject`.

Short config examples:

```yaml
- type: rate_limit
  config:
    limit: 50
    window: 1m

- type: timeout
  config:
    duration: 2s

- type: inject
  config:
    context:
      service_name: gateway
      version: v1
```

`recovery`, `logging`, `request_id`, `real_ip`, and `metric` do not require config.
