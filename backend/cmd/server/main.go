package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lesslyzuniga/calculator-app/backend/internal/handler"
)

func main() {
	mux := newMux()

	log.Println("Calculator API running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/add", handler.Add)
	mux.HandleFunc("/api/v1/divide", handler.Divide)
	mux.HandleFunc("/api/v1/multiply", handler.Multiply)
	mux.HandleFunc("/api/v1/percentage", handler.Percentage)
	mux.HandleFunc("/api/v1/power", handler.Power)
	mux.HandleFunc("/api/v1/square-root", handler.SquareRoot)
	mux.HandleFunc("/api/v1/subtract", handler.Subtract)
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]string{
					"code":    "METHOD_NOT_ALLOWED",
					"message": "Only GET requests are allowed",
				},
			}); err != nil {
				log.Printf("failed to encode health error response: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		}); err != nil {
			log.Printf("failed to encode health response: %v", err)
		}
	})

	return mux
}
