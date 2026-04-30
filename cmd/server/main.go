package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
