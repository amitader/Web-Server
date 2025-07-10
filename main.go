package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	newServer := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	newServer.ListenAndServe()
}