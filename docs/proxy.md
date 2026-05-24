# Описание Proxy

Краткое описание работы слоя proxy.

В проекте за основу взят [ReverseProxy](https://pkg.go.dev/net/http/httputil#ReverseProxy) из пакета [httputil](https://pkg.go.dev/net/http/httputil). С кастомным `Rewrite` свойством. 