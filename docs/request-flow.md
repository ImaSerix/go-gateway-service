# Request Flow

1. Global middleware from `server.middlewares`.
2. Routing by `routes[].path` and `routes[].method`.
3. Local middleware from `routes[].middlewares`.
4. Checks from `routes[].checks`.
5. Transforms from `routes[].transforms`.
6. Proxying to `routes[].upstream`.

A `policy` check can perform a separate internal request during the checks step. `store` only makes sense after such a response exists, because Store reads `http.Response`, not the incoming `http.Request`.
