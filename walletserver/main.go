package main

import (
	"net/http"

	"synnergy/internal/log"
)

func main() {
	srv := newServer()
	http.HandleFunc("/health", srv.healthHandler)
	http.HandleFunc("/wallet/new", srv.newWalletHandler)

	addr := ":8080"
	log.Info("wallet server listening", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("server shutdown", "err", err)
	}
}
