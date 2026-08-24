package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/lesslyzuniga/calculator-app/backend/internal/calculator"
)

const invalidRequestMessage = "Exactly two numeric operands are required"

type addRequest struct {
	Operands []float64 `json:"operands"`
}

type addResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{
			Error: apiError{
				Code:    "METHOD_NOT_ALLOWED",
				Message: "Only POST requests are allowed",
			},
		})
		return
	}

	var request addRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeInvalidRequest(w)
		return
	}

	if err := ensureEndOfBody(decoder); err != nil {
		writeInvalidRequest(w)
		return
	}

	result, err := calculator.Add(request.Operands)
	if err != nil {
		writeInvalidRequest(w)
		return
	}

	writeJSON(w, http.StatusOK, addResponse{Result: result})
}

func ensureEndOfBody(decoder *json.Decoder) error {
	var value any
	if err := decoder.Decode(&value); err != io.EOF {
		return errors.New("request body must contain a single JSON value")
	}

	return nil
}

func writeInvalidRequest(w http.ResponseWriter) {
	writeJSON(w, http.StatusBadRequest, errorResponse{
		Error: apiError{
			Code:    "INVALID_REQUEST",
			Message: invalidRequestMessage,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
