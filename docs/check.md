# Checks

A check runs after route middleware and before route transforms/proxying. It can either allow the request to continue or return an error.

Supported checks:
- `policy`
- `header_required`
- `ip_whitelist`
- `query_required`

`auth`, `inject`, `rate_limit`, and `timeout` are no longer supported as checks. `auth` was replaced by the more general `policy` check, `inject` exists as middleware, and rate limit/timeout behavior also lives in middleware.

## policy

`policy` performs an internal HTTP request, validates the response status code, and can store response values in context on success.

```yaml
checks:
  - type: policy
    config:
      transform:
        header:
          X-Request-ID: "{header:X-Request-ID}"
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

Store inside `policy` works only with the response:
- `{header:key}` reads a response header.
- `{body:key}` reads a top-level JSON body field.

Complex JSON paths such as `{body:user.id}` are not supported yet.

## header_required

Checks that required headers are present in the incoming request.

```yaml
checks:
  - type: header_required
    config:
      headers:
        - X-Request-ID
        - Authorization
```

## ip_whitelist

Checks the IP from `RemoteAddr` without the port.

```yaml
checks:
  - type: ip_whitelist
    config:
      ips:
        - 127.0.0.1
```

## query_required

Checks that required query parameters are present.

```yaml
checks:
  - type: query_required
    config:
      query_params:
        - locale
        - page
```
