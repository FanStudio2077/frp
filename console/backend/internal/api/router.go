package api

import (
	"encoding/json"
	"net/http"
)

// NewRouter builds the console HTTP router skeleton.
// Resource routes (auth/users/servers/nodes/instances/proxies/visitors/overview/audit/alerts)
// will be registered here as the console is implemented.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	return mux
}
