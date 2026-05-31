# Файл конфигурации сервисса `config.yaml`

Эта документация описывает структуру `config.yaml`.

## Root объект

```yaml
server: {}
routes: []
```
| Field             | Type      | Description                       |
| ---               | ---       | ---                               |
| [server](#server) | объект    | глобальная конфигурация сервера   |
| [routes](#route) | список      | список определений путей сервиса  |

## Server
Глобальная конфигурация сервера.

#### Структура 

```yaml
server:
    middlewares: []

```

### Поля

| Field                     | Type      | Description                       |
| ---                       | ---       | ---                               |
| [middlewares](#middleware) | список    | глобально применяемые middleware, выполняются по порядку   |

#### Пример

```yaml
server:
    middleawares:
        - type: rate_limit
          config:
            limit: 50
            window: 1m
        - type: logger
```

### Route
Список определений роутов сервиса.

#### Структура 

```yaml
    routes:
        - path: `строка`
          method: `строка`
          middlewares: [] (опциональное)
          checks: [] (опциональное)
          transforms: {} (опциональное)
          upstream: {}
```

#### Поля

| Field | Type | Description |
| --- | --- | --- |
| path | строка | путь по которому будет доступен роут, может содержать шаблонные значения, например, `/user/{id}`
| method | строка | метод по которому будет доступен роут
| [middlewares](#Middleware) | список | применяемые к роуту middleware
| [checks](#Check) | список | применяемые к роуту проверки
| [transforms](#Transform) | объект | преобразования начального запроса
| [upstream](#Upstream) | объект | конфигурация проксирования

#### Пример

```yaml

route:
    - path: /chat/{id}
      method: GET
      middlewares:
        - type: timeout
          config:
            duration: 2s
      checks:
        - type: header_required
          config:
            headers:
                - X-Username
                - X-Password  
        - type: auth
          config:
            url: "http://auth.url"
            headers:
                X-Username: X-Username
                X-Password: X-Password
            method: POST
            store:
                headers:
                    token: X-Token
            expected_status: 200
     transforms:
        headers:
            X-Token: {context:token}
        query: 
            limit: {query:limit} 
     upstream:
        host: http://example.host
        scheme: http
        path: /user/{route:id}
        method: GET

    - path: /ping
      method: GET
      upstream:
        host: http://example.host
        scheme: http
        path: /ping
        method: GET

```

### Middleware
Описание и настройка миддлеваре.

### Структура

```yaml
middleware:
  - type: string
    config: {}
```

### Поля

| Field | Type | Description |
| --- | --- | --- |
| type | string | Тип middleware
| config | object | Настройки middleware. Некоторые middleware не требуют настройки.

### Типы и небольшие сниппеты настроек

#### CORS

```yaml
middleware:
  - type: cors
    config: 
      origin:
        - http://example.host
        - 127.0.0.1
      method:
        - GET
        - POST
      header:
        - Content-Type
        - X-Request-ID 
        - Authorization
``` 

Все тэги config обязательны. Если какого-то тэга не будет, то этот тэг не будет отдаваться при CORS запросах, что может привести к неожидаемому поведению.

#### Recovery

```yaml
middleware:
  - type: recovery
```

#### RateLimit

```yaml
middleware:
  - type: rate_limit
    config:
      limit: int
      window: duration
```

- limit: любое позитивное число. Кол-во запросов на клиента за момент времени размером window.
- window: время сброса ограничения. Валидный time.Duration (например, 1s, 2m, 1h и другие).

#### Logging

```yaml
middleware:
  - type: logging
```

#### RequestID

```yaml
middleware:
  - type: request_id
```

#### RealIP

```yaml
middleware:
  - type: real_ip
```

#### TimeOut

```yaml
middleware:
  - type: timeout
    config:
      duration: duration
```

- duration: время на выполнение запроса. Валидный time.Duration (например, 1s, 2m, 1h и другие).

#### Metric

```yaml
middleware:
  - type: metric
```

#### Inject

```yaml
middleware:
  - type: inject
    config:
      context:
        name: {query:name}
```
- context: ключи и значения, которые надо добавить в контекст

### Check
Описание и настройка чеков.

#### Policy

```yaml
check:
  - type: policy
    config:
      transform: (Transform config object)
        header:
          X-Username: {header:X-Username}
          X-Password: {header:X-Password}
      upstream: (Upstream config object)
          host: example.host
          scheme: http
          path: /auth (optional)
          method: POST (default, POST)
      expected_status: 200
      store:  (Store config object)
        token: {header:X-Token}
```

- [transform](#transform)
- [upstream](#upstream)
- expected_status: ожидаемый статус ответа usptream
- [store](#store)

#### Header required

```yaml
check:
  - type: header_required
    config:
      header:
        - X-Username
        - X-Password
```

- header: ожидаемые хэдэры в запросе

#### IPWhiteList

```yaml
check:
  - type: ip_whitelist
    config:
      ip:
        - 127.0.0.1
        - 192.168.0.1
```

- ip: разрешённые ip, без порта

#### Query required

```yaml
check:
  - type: query_required
    config:
      query:
        - age
        - name
```

- query: ожидаемые GET параметры

### Store
Описание сохраняемых в контекст пар ключ - значение.

Может содержать произвольные ключ-значение.

### Структура

```yaml
store: {}
```

### Пример

```yaml
store:
  token: {header:X-Token}
``` 

### Transform
Описание преображений запроса.

#### Структура

```yaml
transform:
  query: {}
  header: {}
  body: {}
```

#### Поля

| Field | Type | Description |
| --- | --- | --- |
| query | объект | Пары, ключ: значение, которое сохраняется в GET параметры запроса
| header | объект | Пары, ключ: значение, по ключу сохраняется значение в хэдэр
| body | объект | Пары, ключ: значение, по ключу сохраняется значение в боди (значение также может быть объектом)

#### Пример

```yaml
transform:
  query:
    language: en
  header:
    X-Token: {ctx:token}
  body:
    user:
      id: {ctx:user_id}
```

### Upstream
Описание запроса, проксируемого сервером.

#### Структура

```yaml
upstream:
  host: string
  scheme: string
  path: string
  method: string
```

#### Поля

| Field | Type | Description |
| --- | --- | --- |
| host | string | Хост к которому обращаются
| scheme | string | Схема по которой обращаются
| path | string | Путь по которому обращаются
| method | string | Метод который используют

#### Пример

```yaml
upstream:
  host: example.host
  scheme: http
  path: /users (optional)
  method: GET (optional)
```

Тэг method и path не обязателен (при их оттуствии, копируется с настроек [route](#route)).  