package main

import (
	"encoding/json"
	"net/http"

	"synnergy/core"
	"synnergy/internal/log"
)

type server struct{}

func newServer() *server { return &server{} }

func (s *server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *server) newWalletHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	wallet, err := core.NewWallet()
	if err != nil {
		log.Error("wallet creation failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"address": wallet.Address})
}
