package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
)

func main() {

	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println(err)
	}

	mux := http.NewServeMux()
	for _, route := range cfg.Routes {
		e, err := endpoint.NewEndpointFromConfig(&route)
		if err != nil {
			slog.Error("failed register endpoint", "path", route.Path, "error", err)
			continue
		}
		mux.Handle(e.Pattern(), e)
	}

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("/pingv2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong pong"))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
