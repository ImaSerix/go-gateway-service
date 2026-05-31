package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder"
	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	clientbuilder "github.com/ImaSerix/go-gateway-service/internal/builder/client"
	"github.com/ImaSerix/go-gateway-service/internal/builder/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/builder/handler"
	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/builder/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/builder/render"
	"github.com/ImaSerix/go-gateway-service/internal/builder/resolver"
	storebuilder "github.com/ImaSerix/go-gateway-service/internal/builder/store"
	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	resolverpkg "github.com/ImaSerix/go-gateway-service/internal/resolver"
)

func main() {

	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := http.DefaultClient

	checkRegisty := check.NewCheckRegistry()
	middlewareRegistry := middleware.NewMiddlewareRegistry()
	transformerRegistry := transformer.NewTransformerRegistry()

	resolver := resolver.NewMultiResolverBuilder().Build()

	render := render.NewBuilder(resolver).Build()
	responseRender := renderer.NewResponseRender(resolverpkg.NewResponseMultiResolver())

	transformerBuilder := transformer.NewBuilder(transformerRegistry)
	clientBuilder := clientbuilder.NewBuilder(client, render)
	storeBuilder := storebuilder.NewBuilder(responseRender)

	builder.RegisterChecks(checkRegisty, transformerBuilder, clientBuilder, storeBuilder)
	builder.RegisterTransformers(transformerRegistry, render)
	builder.RegisterMiddlewares(middlewareRegistry)

	checkBuilder := check.NewBuilder(checkRegisty)
	middlewareBuilder := middleware.NewBuilder(middlewareRegistry)
	proxyBuilder := proxy.NewBuilder(client, render)
	endpointBuilder := endpoint.NewBuilder(checkBuilder, transformerBuilder, middlewareBuilder, proxyBuilder)

	handlerBuilder := handler.NewBuilder(middlewareBuilder, endpointBuilder)

	h, err := handlerBuilder.Build(cfg)
	if err != nil {
		slog.Error("failed to build handler", "err", err)
		return
	}

	log.Print(http.ListenAndServe(":8080", h))
}
