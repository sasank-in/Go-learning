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
	Expression string             `json:"expression"`
	Operation  string             `json:"operation"`
	A          float64            `json:"a"`
	B          float64            `json:"b"`
	Variables  map[string]float64 `json:"variables,omitempty"`
	// AngleMode is "deg" or "rad" (default). Controls trig functions.
	AngleMode string `json:"angleMode,omitempty"`
}

// CalcResponse is the JSON response returned on success. Variables echoes the
// (possibly updated) variable environment so a stateless client can persist it
// and send it back on the next request.
type CalcResponse struct {
	Result    float64            `json:"result"`
	Variables map[string]float64 `json:"variables,omitempty"`
}

// ErrorResponse is the JSON response returned on failure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register wires the calculator routes onto the given mux, wrapped with
// permissive CORS so the browser frontend can call the API directly.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /calculate", calculate)
	mux.HandleFunc("OPTIONS /calculate", preflight)
}

// withCORS sets permissive CORS headers on the response.
func withCORS(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Content-Type")
}

// preflight answers CORS preflight OPTIONS requests.
func preflight(w http.ResponseWriter, _ *http.Request) {
	withCORS(w)
	w.WriteHeader(http.StatusNoContent)
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
	var vars map[string]float64
	var err error
	angle := calculator.Radians
	if strings.EqualFold(req.AngleMode, "deg") {
		angle = calculator.Degrees
	}
	switch {
	case strings.TrimSpace(req.Expression) != "":
		// Expression form supports variables and assignment; the updated
		// environment is echoed back to the client.
		result, vars, err = calculator.EvaluateWithOptions(req.Expression, calculator.Options{
			Vars:  req.Variables,
			Angle: angle,
		})
	case req.Operation != "":
		// Structured two-operand form.
		result, err = calculator.Compute(req.Operation, req.A, req.B)
	default:
		// Neither form supplied — treat as an empty expression rather than
		// falling through to an opaque "unknown operation".
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "empty expression"})
		return
	}
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, calculator.ErrUnknownOperation):
			status = http.StatusUnprocessableEntity
		case errors.Is(err, calculator.ErrNotFinite):
			status = http.StatusUnprocessableEntity
		}
		writeJSON(w, status, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, CalcResponse{Result: result, Variables: vars})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	withCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
