package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lesslyzuniga/calculator-app/backend/internal/calculator"
)

type multiplyRequest struct {
	Operands []float64 `json:"operands"`
}

type multiplyResponse struct {
	Result float64 `json:"result"`
}

func Multiply(w http.ResponseWriter, r *http.Request) {
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

	var request multiplyRequest
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

	result, err := calculator.Multiply(request.Operands)
	if errors.Is(err, calculator.ErrInvalidResult) {
		writeInvalidResult(w)
		return
	}

	if err != nil {
		writeInvalidRequest(w)
		return
	}

	writeJSON(w, http.StatusOK, multiplyResponse{Result: result})
}
