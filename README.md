# go-gateway-service

Учебный API Gateway на Go, который регистрирует маршруты из YAML-конфига, выполняет middleware и checks, трансформирует запрос и проксирует его в upstream.

Подробные разделы лежат в `./docs`:
- [config.md](./docs/config.md) - полный формат `config.yaml`.
- [check.md](./docs/check.md) - актуальные checks и их ответственность.
- [middleware.md](./docs/middleware.md) - актуальные middleware.
- [transformer.md](./docs/transformer.md) - трансформации запроса.
- [proxy.md](./docs/proxy.md) - слой проксирования.
- [request-flow.md](./docs/request-flow.md) - порядок выполнения запроса.

## Быстрый старт

```bash
go run ./cmd/server -config ./config.yaml
```

Сервер читает `server.middlewares` и `routes`, после чего поднимает HTTP-handler на `:8080`.

## Кратко о конфиге

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

## Актуальные checks

Сейчас зарегистрированы только:
- `policy` - делает внутренний HTTP-запрос, проверяет статус и может сохранить данные из response в context.
- `header_required` - требует наличие headers.
- `ip_whitelist` - пропускает только разрешенные IP.
- `query_required` - требует наличие query-параметров.

Старый `auth` check заменен на `policy` и не считается актуальным.

## Resolver

Шаблоны пишутся в формате `{source:key}`.

Request-renderer используется в proxy, client и transformers:
- `{context:key}` - значение из `request.Context()`.
- `{route:key}` - параметр роутинга chi.
- `{query:key}` - query-параметр входящего запроса.
- `{header:key}` - header входящего запроса.

Response-renderer используется только в `Store`, то есть только там, где код уже сделал внутренний запрос и получил `http.Response`.
- `{header:key}` - header response.
- `{body:key}` - поле верхнего уровня JSON body response.

Ограничение Store намеренное: `body` сейчас читает только верхний уровень JSON-объекта. Вложенные пути, массивы и сложные выражения пока не поддерживаются.

## Store

`Store` сохраняет пары `contextKey: template` в `request.Context()` после внутреннего response. Практически это нужно для checks или middleware, которые сами вызывают внешний сервис и хотят передать результат дальше по цепочке. Текущий встроенный `inject` middleware только кладет литеральные значения в context и не использует Store.

```yaml
store:
  token: "{header:X-Token}"
  user_id: "{body:user_id}"
```

После `policy` эти значения доступны как `{context:token}` и `{context:user_id}` в последующих transforms.
