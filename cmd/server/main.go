package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func main() {

	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	b := builder.NewEndpointBuilder(http.DefaultClient)

	mux := http.NewServeMux()
	for _, route := range cfg.Routes {
		e, err := b.BuildEndpoint(&route)
		if err != nil {
			slog.Error("failed register endpoint", "path", route.Path, "error", err)
			continue
		}
		mux.Handle(e.Pattern(), e)
	}

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Print("i am here")
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("/pingv2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong pong"))
	})

	log.Print(http.ListenAndServe(":8080", mux))
}
