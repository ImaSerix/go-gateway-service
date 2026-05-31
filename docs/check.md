# Checks

Check выполняется после route middleware и до route transforms/proxy. Он может пропустить запрос дальше или вернуть ошибку.

Актуальные checks:
- `policy`
- `header_required`
- `ip_whitelist`
- `query_required`

`auth`, `inject`, `rate_limit` и `timeout` как checks больше не актуальны. `auth` заменен на более общий `policy`, `inject` живет как middleware, а rate limit/timeout живут как middleware.

## policy

`policy` делает внутренний HTTP-запрос, проверяет status code и при успехе может сохранить данные из response в context.

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

Store внутри `policy` работает только с response:
- `{header:key}` читает response header.
- `{body:key}` читает поле верхнего уровня JSON body.

Сложные JSON-пути вроде `{body:user.id}` пока не поддерживаются.

## header_required

Проверяет наличие headers во входящем запросе.

```yaml
checks:
  - type: header_required
    config:
      headers:
        - X-Request-ID
        - Authorization
```

## ip_whitelist

Проверяет IP из `RemoteAddr` без порта.

```yaml
checks:
  - type: ip_whitelist
    config:
      ips:
        - 127.0.0.1
```

## query_required

Проверяет наличие query-параметров.

```yaml
checks:
  - type: query_required
    config:
      query_params:
        - locale
        - page
```
