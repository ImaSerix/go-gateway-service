# Proxy

The proxy layer is based on [`httputil.ReverseProxy`](https://pkg.go.dev/net/http/httputil#ReverseProxy) with a custom `Rewrite` function.

The route `upstream` config defines the target URL and method. The upstream path can contain request templates such as `{route:id}` or `{query:locale}`.
