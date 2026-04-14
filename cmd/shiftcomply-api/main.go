// Command shiftcomply-api is an HTTP server exposing the Shift Comply
// regulation engine as a REST API.
//
// Endpoints:
//
//	GET  /jurisdictions              List all jurisdictions
//	GET  /rules?jurisdiction=X       Query rules (filters: staff, unit, scope, category)
//	GET  /constraints?jurisdiction=X Generate optimizer-ready constraints
//	GET  /compare?left=X&right=Y    Compare two jurisdictions
//	POST /validate                   Validate a schedule against jurisdiction rules
//	GET  /export/:code               Full JSON export of a jurisdiction
//
// Start:
//
//	shiftcomply-api                  # listens on :8080
//	shiftcomply-api -addr :3000     # custom port
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func main() {
	addr := ":8080"
	if len(os.Args) > 2 && os.Args[1] == "-addr" {
		addr = os.Args[2]
	}
	if env := os.Getenv("PORT"); env != "" {
		addr = ":" + env
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/jurisdictions", handleJurisdictions)
	mux.HandleFunc("/rules", handleRules)
	mux.HandleFunc("/constraints", handleConstraints)
	mux.HandleFunc("/compare", handleCompare)
	mux.HandleFunc("/validate", handleValidate)
	mux.HandleFunc("/export/", handleExport)
	mux.HandleFunc("/health", handleHealth)

	handler := withCORS(mux)
	log.Printf("shift-comply API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]any{"status": "ok", "jurisdictions": len(comply.All())})
}

func handleJurisdictions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, comply.All())
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := comply.Code(r.URL.Query().Get("jurisdiction"))
	if code == "" {
		http.Error(w, `{"error":"jurisdiction parameter required"}`, http.StatusBadRequest)
		return
	}
	if comply.For(code) == nil {
		http.Error(w, fmt.Sprintf(`{"error":"unknown jurisdiction: %s"}`, code), http.StatusNotFound)
		return
	}

	opts := buildQueryOpts(r)
	writeJSON(w, comply.EffectiveRules(code, opts...))
}

func handleConstraints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := comply.Code(r.URL.Query().Get("jurisdiction"))
	if code == "" {
		http.Error(w, `{"error":"jurisdiction parameter required"}`, http.StatusBadRequest)
		return
	}
	if comply.For(code) == nil {
		http.Error(w, fmt.Sprintf(`{"error":"unknown jurisdiction: %s"}`, code), http.StatusNotFound)
		return
	}

	opts := buildQueryOpts(r)
	writeJSON(w, comply.GenerateConstraints(code, opts...))
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	left := comply.Code(r.URL.Query().Get("left"))
	right := comply.Code(r.URL.Query().Get("right"))
	if left == "" || right == "" {
		http.Error(w, `{"error":"left and right parameters required"}`, http.StatusBadRequest)
		return
	}

	opts := buildQueryOpts(r)
	writeJSON(w, comply.Compare(left, right, opts...))
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var schedule comply.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"invalid JSON: %s"}`, err), http.StatusBadRequest)
		return
	}

	report, err := comply.Validate(schedule)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusBadRequest)
		return
	}

	writeJSON(w, report)
}

func handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := comply.Code(strings.TrimPrefix(r.URL.Path, "/export/"))
	if code == "" {
		http.Error(w, `{"error":"jurisdiction code required in path"}`, http.StatusBadRequest)
		return
	}

	j := comply.For(code)
	if j == nil {
		http.Error(w, fmt.Sprintf(`{"error":"unknown jurisdiction: %s"}`, code), http.StatusNotFound)
		return
	}

	writeJSON(w, j)
}

func buildQueryOpts(r *http.Request) []comply.QueryOption {
	var opts []comply.QueryOption
	if s := r.URL.Query().Get("staff"); s != "" {
		opts = append(opts, comply.ForStaff(comply.Key(s)))
	}
	if s := r.URL.Query().Get("unit"); s != "" {
		opts = append(opts, comply.ForUnit(comply.Key(s)))
	}
	if s := r.URL.Query().Get("scope"); s != "" {
		opts = append(opts, comply.ForScope(comply.Scope(s)))
	}
	if s := r.URL.Query().Get("category"); s != "" {
		opts = append(opts, comply.InCategory(comply.Category(s)))
	}
	return opts
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
