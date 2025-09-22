package main

import (
	"net/http"

	"synnergy/internal/log"
)

var httpListenAndServe = http.ListenAndServe

func run(addr string, srv *server) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.healthHandler)
	mux.HandleFunc("/wallet/new", srv.newWalletHandler)
	log.Info("wallet server listening", "addr", addr)
	return httpListenAndServe(addr, mux)
}

func main() {
	if err := run(":8080", newServer()); err != nil {
		log.Error("server shutdown", "err", err)
	}
}
