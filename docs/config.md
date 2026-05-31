# Конфигурация

Файл `config.yaml` состоит из `server` и `routes`.

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

`server.middlewares` - глобальные middleware, применяются ко всем маршрутам.

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

Поля:
- `path` - входящий путь. Может содержать chi-параметры, например `/users/{id}`.
- `method` - входящий HTTP-метод.
- `middlewares` - middleware только для этого route.
- `checks` - проверки перед transform/proxy.
- `transforms` - изменения запроса перед proxy.
- `upstream` - целевой сервис.

## Upstream

```yaml
upstream:
  scheme: http
  host: api.local
  path: /users/{route:id}
  method: GET
```

Поля:
- `scheme` - `http`, `https` или `ws`.
- `host` - host без scheme.
- `path` - путь upstream, может содержать request-шаблоны.
- `method` - метод upstream. Для route upstream при пустом значении берется `route.method`.

## Resolver

Шаблон: `{source:key}`.

Для request-renderer доступны:
- `{context:key}`
- `{route:key}`
- `{query:key}`
- `{header:key}`

Для response-renderer в `store` доступны только:
- `{header:key}`
- `{body:key}`

`{body:key}` читает только поле верхнего уровня JSON-объекта response body. Вложенные пути и массивы пока не поддерживаются.

## Checks

Актуальные checks соответствуют регистрации в `internal/builder/bootstrap.go`: `policy`, `header_required`, `ip_whitelist`, `query_required`.

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

`policy` делает внутренний запрос в `upstream`, применяя к нему `transform`, сравнивает response status с `expected_status`, а затем сохраняет данные через `store`.

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

Store работает с `http.Response`, а не с входящим `http.Request`. Поэтому он уместен только в checks или middleware, которые делают внутренний запрос и получают response.

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

Типы:
- `header` - выставляет headers.
- `query_params` - выставляет query-параметры.
- `body_fields` - мержит поля в JSON body. Сейчас рекурсивно поддерживаются объекты, но не массивы.

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

Актуальные middleware: `cors`, `recovery`, `rate_limit`, `logging`, `request_id`, `real_ip`, `timeout`, `metric`, `inject`.

Короткие конфиги:

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

`recovery`, `logging`, `request_id`, `real_ip` и `metric` не требуют config.
