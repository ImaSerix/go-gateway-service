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

The internal request is a copy of the client request (method and URL come from `upstream`, headers are copied, then `transform` is applied). The client body is sent only when `forward_body` is `true`. `upstream.method` is `GET` or `POST`.

Redirects are not followed: the status of the first response is the one that is checked, so a `3xx` answer is an unexpected status.

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
      forward_body: false
      store:
        token: "{header:X-Token}"
        user_id: "{body:user_id}"
```

Fields:
- `expected_status` - status that lets the request through. Defaults to `200`.
- `forward_body` - whether the client body is sent with the internal request. Defaults to `false`: the internal request has no body and the client body is not read. With `true` the body is read into memory and sent to both the policy service and the upstream.

Outcome for the client:

| Policy service | Client gets |
| --- | --- |
| `expected_status` | request continues |
| `401` | `401`, code `unauthorized` |
| unreachable, or any other status (a redirect included) | `503`, code `auth_unavailable` |
| `expected_status`, but `store` cannot read a value from the response | `503`, code `auth_unavailable` |
| not called: `transform` of the internal request failed | `500`, code `internal_error` |
| not called: a value from the request turns a path segment of `upstream.path` into `.` or `..` ([pattern.md](pattern.md#values-in-a-url)) | `400`, code `bad_request` |

The body of the policy response is never passed to the client. Error format: [request-flow.md](request-flow.md#gateway-errors).

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
