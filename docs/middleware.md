# Middleware

Middleware может быть глобальным (`server.middlewares`) или локальным (`routes[].middlewares`). Локальные middleware выполняются после роутинга и до checks.

Актуальные middleware:
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

`allowed.origin` обязателен. `method` и `header` можно оставить пустыми, но тогда CORS-ответ будет без этих разрешений.

## recovery

```yaml
middlewares:
  - type: recovery
```

Перехватывает panic и возвращает HTTP 500.

## rate_limit

```yaml
middlewares:
  - type: rate_limit
    config:
      limit: 50
      window: 1m
```

Ограничивает количество запросов за окно `window`.

## logging

```yaml
middlewares:
  - type: logging
```

Логирует базовую информацию по запросу.

## request_id

```yaml
middlewares:
  - type: request_id
```

Добавляет request id в context/headers.

## real_ip

```yaml
middlewares:
  - type: real_ip
```

Определяет клиентский IP по proxy headers, если они есть.

## timeout

```yaml
middlewares:
  - type: timeout
    config:
      duration: 2s
```

Ограничивает время обработки запроса.

## metric

```yaml
middlewares:
  - type: metric
```

Собирает простые метрики запросов.

## inject

```yaml
middlewares:
  - type: inject
    config:
      context:
        service_name: gateway
        version: v1
```

Добавляет литеральные значения в `request.Context()`. Сейчас этот middleware не делает внутренний запрос и не использует Store.
