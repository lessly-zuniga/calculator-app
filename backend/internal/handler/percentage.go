package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lesslyzuniga/calculator-app/backend/internal/calculator"
)

type percentageRequest struct {
	Operands []float64 `json:"operands"`
}

type percentageResponse struct {
	Result float64 `json:"result"`
}

func Percentage(w http.ResponseWriter, r *http.Request) {
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

	var request percentageRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writePercentageInvalidRequest(w)
		return
	}

	if err := ensureEndOfBody(decoder); err != nil {
		writePercentageInvalidRequest(w)
		return
	}

	result, err := calculator.Percentage(request.Operands)
	if err != nil {
		writePercentageInvalidRequest(w)
		return
	}

	writeJSON(w, http.StatusOK, percentageResponse{Result: result})
}

func writePercentageInvalidRequest(w http.ResponseWriter) {
	writeJSON(w, http.StatusBadRequest, errorResponse{
		Error: apiError{
			Code:    "INVALID_REQUEST",
			Message: "Exactly one numeric operand is required",
		},
	})
}
