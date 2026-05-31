# Порядок выполнения запроса

1. Глобальные middleware из `server.middlewares`.
2. Роутинг по `routes[].path` и `routes[].method`.
3. Локальные middleware из `routes[].middlewares`.
4. Checks из `routes[].checks`.
5. Transforms из `routes[].transforms`.
6. Proxy в `routes[].upstream`.

`policy` check внутри шага checks может выполнить отдельный внутренний запрос. Только после такого response имеет смысл использовать `store`, потому что Store читает именно `http.Response`, а не входящий `http.Request`.
