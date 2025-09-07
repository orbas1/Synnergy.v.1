package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	synn "synnergy"
)

// newHandler exposes a /metrics endpoint returning watchtower health metrics in
// JSON form. It allows external systems to poll the node without embedding CLI
// logic.
func newHandler(wt *synn.WatchtowerNode) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		m := wt.Metrics()
		b, _ := json.Marshal(m)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	})
	return mux
}

func main() {
	wt := synn.NewWatchtowerNode("monitor", nil)
	if err := wt.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{Addr: ":9090", Handler: newHandler(wt)}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
