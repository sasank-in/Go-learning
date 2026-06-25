// Package handlers exposes the calculator over HTTP as a small JSON API.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"calculator-application/backend/internal/calculator"
)

// CalcRequest is the expected JSON request body for /calculate.
//
// Two forms are supported:
//   - Expression form: {"expression": "2 + 3 * (4 - 1)"} — evaluates any
//     arithmetic expression.
//   - Structured form: {"operation": "add", "a": 2, "b": 3} — a single
//     two-operand operation.
//
// If "expression" is non-empty it takes precedence.
type CalcRequest struct {
	Expression string  `json:"expression"`
	Operation  string  `json:"operation"`
	A          float64 `json:"a"`
	B          float64 `json:"b"`
}

// CalcResponse is the JSON response returned on success.
type CalcResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse is the JSON response returned on failure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register wires the calculator routes onto the given mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /calculate", calculate)
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func calculate(w http.ResponseWriter, r *http.Request) {
	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	var result float64
	var err error
	if strings.TrimSpace(req.Expression) != "" {
		result, err = calculator.Evaluate(req.Expression)
	} else {
		result, err = calculator.Compute(req.Operation, req.A, req.B)
	}
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, calculator.ErrUnknownOperation) {
			status = http.StatusUnprocessableEntity
		}
		writeJSON(w, status, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, CalcResponse{Result: result})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
