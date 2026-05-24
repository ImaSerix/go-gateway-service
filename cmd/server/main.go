package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder"
	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	"github.com/ImaSerix/go-gateway-service/internal/builder/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/builder/handler"
	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/builder/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/builder/render"
	"github.com/ImaSerix/go-gateway-service/internal/builder/resolver"
	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

//TODO Может быть имеет смысл сделать ExternalPolicyCheck, что-то похожее на auth но немного другое. И вероятно сделать какой-нибудь универсальный чек, и на auth и на External

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

	builder.RegisterChecks(checkRegisty, render, client)
	builder.RegisterTransformers(transformerRegistry, render)
	builder.RegisterMiddlewares(middlewareRegistry)

	checkBuilder := check.NewBuilder(checkRegisty)
	middlewareBuilder := middleware.NewBuilder(middlewareRegistry)
	transformerBuilder := transformer.NewBuilder(transformerRegistry)
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
