# Proxy

The proxy layer is based on [`httputil.ReverseProxy`](https://pkg.go.dev/net/http/httputil#ReverseProxy) with a custom `Rewrite` function.

The route `upstream` config defines the target URL and method. The upstream path can contain request templates such as `{route:id}` or `{query:locale}`.

## Upstream path

`upstream.path` decides which path the upstream receives:

| `upstream.path` | Incoming request | Upstream request |
| --- | --- | --- |
| `/resource/{route:id}` | `/resource/42?limit=5` | `/resource/42?limit=5` |
| `/api/v1/profiles/{route:id}` | `/users/42` | `/api/v1/profiles/42` |
| empty | `/resource/42?limit=5` | `/resource/42?limit=5` |
| `/search?source=gateway` | `/find?q=cats` | `/search?source=gateway&q=cats` |

- A non-empty path replaces the incoming path. Templates are rendered first.
- An empty path passes the incoming path through unchanged.
- The client query string is always passed on. If `upstream.path` has its own query, the parameters are joined: first the ones from the config, then the ones from the client. A parameter present in both is sent twice.

## Upstream URL errors

If a template in `upstream.path` cannot be rendered (for example `{query:locale}` without `locale` in the request) or the rendered URL cannot be parsed, the gateway answers `500` with code `internal_error` and writes the reason to its log. The request is not sent to the upstream. To reject requests without a required parameter with `403`, use the `query_required` or `header_required` check.

`Rewrite` cannot return an error, so it cancels the outgoing request with a cause and clears its address; `ErrorHandler` recognizes the cause and writes the `500`. Any other transport error is answered with `502` `upstream_unavailable` ([request-flow.md](request-flow.md#gateway-errors)).
