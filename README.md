# go-gateway-service

A learning API Gateway in Go that registers routes from YAML config, runs middleware and checks, transforms requests, and proxies them to upstream services.

More detailed documentation lives in `./docs`:
- [config.md](./docs/config.md) - full `config.yaml` format.
- [check.md](./docs/check.md) - supported checks and their responsibilities.
- [middleware.md](./docs/middleware.md) - supported middleware.
- [transformer.md](./docs/transformer.md) - request transforms.
- [proxy.md](./docs/proxy.md) - proxy layer.
- [request-flow.md](./docs/request-flow.md) - request execution order.

## Quick Start

```bash
go run ./cmd/server -config ./config.yaml
```

The server reads `server.middlewares` and `routes`, then starts an HTTP handler on `:8080`.

## Config Overview

```yaml
server:
  middlewares:
    - type: recovery
    - type: logging

routes:
  - path: /users/{id}
    method: GET
    middlewares:
      - type: request_id
    checks:
      - type: header_required
        config:
          headers:
            - X-Request-ID
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
    transforms:
      header:
        Authorization: "Bearer {context:token}"
    upstream:
      scheme: http
      host: api.local
      path: /users/{route:id}
      method: GET
```

## Supported Checks

Currently registered checks:
- `policy` - performs an internal HTTP request, validates the response status, and can store response values in context.
- `header_required` - requires incoming request headers.
- `ip_whitelist` - allows only configured IP addresses.
- `query_required` - requires incoming query parameters.

The old `auth` check has been replaced by `policy` and is no longer supported.

## Resolver

Templates use the `{source:key}` format.

The request renderer is used by proxy, client, and transformers:
- `{context:key}` - value from `request.Context()`.
- `{route:key}` - chi route parameter.
- `{query:key}` - incoming query parameter.
- `{header:key}` - incoming request header.

The response renderer is used only by `Store`, which means only code that already made an internal request and received an `http.Response` can use it:
- `{header:key}` - response header.
- `{body:key}` - top-level JSON body field from the response.

The Store limitation is intentional: `body` currently reads only top-level fields from a JSON object. Nested paths, arrays, and complex expressions are not supported yet.

## Store

`Store` saves `contextKey: template` pairs into `request.Context()` after an internal response is available. In practice, this is useful for checks or middleware that call an external service and need to pass its response data further down the request chain. The built-in `inject` middleware only stores literal values in context and does not use Store.

```yaml
store:
  token: "{header:X-Token}"
  user_id: "{body:user_id}"
```

After `policy` runs, these values are available to later transforms as `{context:token}` and `{context:user_id}`.
