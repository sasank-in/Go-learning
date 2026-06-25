// Command server starts the calculator HTTP API.
package main

import (
	"log"
	"net/http"
	"os"

	"calculator-application/backend/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	handlers.Register(mux)

	addr := ":" + port
	log.Printf("calculator backend listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
