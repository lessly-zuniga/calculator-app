package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lesslyzuniga/calculator-app/backend/internal/calculator"
)

type squareRootRequest struct {
	Operands []float64 `json:"operands"`
}

type squareRootResponse struct {
	Result float64 `json:"result"`
}

func SquareRoot(w http.ResponseWriter, r *http.Request) {
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

	var request squareRootRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeSquareRootInvalidRequest(w)
		return
	}

	if err := ensureEndOfBody(decoder); err != nil {
		writeSquareRootInvalidRequest(w)
		return
	}

	result, err := calculator.SquareRoot(request.Operands)
	if errors.Is(err, calculator.ErrNegativeSquareRoot) {
		writeJSON(w, http.StatusUnprocessableEntity, errorResponse{
			Error: apiError{
				Code:    "NEGATIVE_SQUARE_ROOT",
				Message: "Cannot calculate the square root of a negative number",
			},
		})
		return
	}

	if errors.Is(err, calculator.ErrInvalidResult) {
		writeInvalidResult(w)
		return
	}

	if err != nil {
		writeSquareRootInvalidRequest(w)
		return
	}

	writeJSON(w, http.StatusOK, squareRootResponse{Result: result})
}

func writeSquareRootInvalidRequest(w http.ResponseWriter) {
	writeJSON(w, http.StatusBadRequest, errorResponse{
		Error: apiError{
			Code:    "INVALID_REQUEST",
			Message: "Exactly one numeric operand is required",
		},
	})
}
