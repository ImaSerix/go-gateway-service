# go-gateway-service

Учебный API Gateway на Go (chi), который проксирует запросы на upstream-сервисы, поддерживает middleware, checks, transform-пайплайн и шаблонный resolver.

> ⚠️ Проект учебный: архитектура и реализация могут быть упрощены, местами намеренно оставлены шероховатости.

## Содержание
- [Быстрый старт](#быстрый-старт)
- [Возможности](#возможности)
- [Конфигурация](#конфигурация)
  - [Корневые ключи](#корневые-ключи)
  - [Описание `route`](#описание-route)
  - [Upstream](#upstream)
  - [Checks](#checks)
  - [Middleware](#middleware)
  - [Transform](#transform)
- [Resolver и шаблоны значений](#resolver-и-шаблоны-значений)
- [Store и контекст](#store-и-контекст)
- [Пример config.yaml](#пример-configyaml)
- [TODO](#todo)

## Быстрый старт
1. Подготовьте `config.yaml`.
2. Запустите:
```bash
go run ./cmd/server -config ./config.yaml
```
3. Сервис поднимет HTTP-сервер и зарегистрирует маршруты из `routes`.

## Возможности
- Глобальные и роут-специфичные middleware.
- Проверки (checks) перед проксированием запроса.
- Преобразование body/header перед отправкой в upstream.
- Подстановка значений из запроса/контекста через resolver.
- Гибкая YAML-конфигурация без перекомпиляции.

## Конфигурация

### Корневые ключи
- `server`: настройки сервера, сейчас в основном список глобальных middleware.
- `routes`: список правил маршрутизации.

### Описание `route`
Для каждого элемента `routes`:
- `path`: путь входящего запроса (например `/users/{id}` или `/posts`).
- `method`: HTTP-метод входящего запроса (`GET`, `POST`, ...).
- `middleware`: middleware, применяемые только к этому route.
- `checks`: проверки, выполняются до proxy.
- `transform`: модификации запроса (header/body).
- `upstream`: куда и как проксировать запрос.

### Upstream
- `url` (string): полный адрес целевого ресурса.
- `method` (string, optional): метод вызова upstream.
  - Если не указан, берётся `route.method`.

### Checks
Общая структура элемента:
```yaml
- type: <check_name>
  config:
    ...
```

Поддерживаемые `type` и ключи `config`:

1. `auth`
   - `url`: URL сервиса авторизации.
   - `method`: HTTP-метод запроса в auth-сервис.
   - `forward_headers`: карта `куда_в_auth: откуда_из_входящего`.
   - `expected_status`: ожидаемый HTTP-код от auth-сервиса.
   - `store`: что сохранить из ответа auth в контекст.
     - `body`: map `ctxKey: jsonField`
     - `headers`: map `ctxKey: headerName`

2. `required_header`
   - `header`: список обязательных headers.

3. `required_query`
   - `query`: список обязательных query-параметров.

4. `ip_whitelist`
   - `ip`: список разрешённых IP.

5. `rate_limit`
   - `limit`: лимит запросов.
   - `window`: окно времени (например `1m`).

6. `inject`
   - `ctx`: map `ctxKey: value`, принудительно добавляет значения в context.

7. `timeout`
   - `duration`: таймаут на check (например `2s`).

### Middleware
Общая структура:
```yaml
- type: <middleware_name>
  config:
    ...
```

Типы middleware:
- `cors`
  - `allowed.origin`: список origin.
  - `allowed.method`: список HTTP-методов.
  - `allowed.header`: список headers.
- `recovery`
- `logging`
- `request_id`
- `real_ip`
- `timeout`
  - `duration`: например `2s`.
- `metric`
- `rate_limit`
  - `limit`, `window`.

### Transform
- `transform.header`: map `headerName: resolver.key`.
- `transform.body`: map/объект шаблона body, где значения могут ссылаться на resolver-ключи.

Пример:
```yaml
transform:
  header:
    X-User-ID: query.user_id
  body:
    authToken: context.token
```

## Resolver и шаблоны значений
Формат ключа: `<source>.<name>`.

Поддерживаемые `source`:
- `query.<key>` — query-параметры.
- `header.<key>` — входящие HTTP headers.
- `router.<key>` — path params (chi URL params).
- `context.<key>` — значения из `request.Context()`.

Если ключ невалидный или source не зарегистрирован, резолвинг вернёт `not found`.

## Store и контекст
Ключевой практический момент:
- В `auth.store` и `inject.ctx` данные сохраняются в `request.Context()` по указанному `ctxKey`.
- Потом эти значения можно использовать в `transform` через `context.<ctxKey>`.

Пример цепочки:
1. `auth` достал `token` и сохранил как `ctx.auth_token` (фактически ключ `auth_token` в context).
2. В `transform.header` указываете:
   - `Authorization: context.auth_token`
3. Перед проксированием header будет заполнен этим значением.

## Пример config.yaml
```yaml
server:
  middleware:
    - type: recovery
    - type: logging
    - type: request_id

routes:
  - path: /users/{id}
    method: GET
    checks:
      - type: required_query
        config:
          query: ["locale"]
      - type: inject
        config:
          ctx:
            service_name: gateway
    transform:
      header:
        X-Request-Locale: query.locale
        X-Service-Name: context.service_name
    upstream:
      url: https://jsonplaceholder.typicode.com/users
      method: GET
```

## TODO
- Привести имена некоторых типов/ключей к единообразному стилю (`required_header` vs `header_required` и т.д.).
- Добавить раздел с ошибками и troubleshooting (как диагностировать проблемы в check/transform).
- Добавить примеры для каждого check/middleware в отдельные snippet-блоки.
- Добавить e2e тест с полным пайплайном: check → store(context) → transform → proxy.
- Вынести и задокументировать matrix совместимости источников resolver с transform/checks.
