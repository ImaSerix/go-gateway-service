# Request Flow

1. Global middleware from `server.middlewares`.
2. Routing by `routes[].path` and `routes[].method`.
3. Local middleware from `routes[].middlewares`.
4. Checks from `routes[].checks`.
5. Transforms from `routes[].transforms`.
6. Proxying to `routes[].upstream`.

A `policy` check can perform a separate internal request during the checks step. `store` only makes sense after such a response exists, because Store reads `http.Response`, not the incoming `http.Request`.

## Gateway errors

When the gateway itself rejects a request, it answers with JSON and `Content-Type: application/json`:

```json
{"code": "unauthorized", "message": "authentication required"}
```

| Situation | Status | `code` |
| --- | --- | --- |
| `policy`: the policy service answered `401` | `401` | `unauthorized` |
| `policy`: the policy service is unreachable, answered with a status other than the expected one and `401` (a redirect included), or `store` cannot read its response | `503` | `auth_unavailable` |
| No route for the path | `404` | `not_found` |
| The path exists, the method does not match | `405` | `method_not_allowed` |
| `header_required`, `ip_whitelist` or `query_required` check failed | `403` | `forbidden` |
| `rate_limit` exceeded | `429` | `too_many_requests` |
| The request to the upstream failed | `502` | `upstream_unavailable` |
| Internal gateway error: panic (`recovery` middleware), failed transform, any other check error | `500` | `internal_error` |

`message` is a fixed text: it has no tokens, internal addresses or internal error text. The reason is written to the gateway log. Responses of upstream services are passed through unchanged.

All of them are written by `httperror.Write` (`internal/httperror`).
