package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesslyzuniga/calculator-app/backend/internal/handler"
)

type apiResponse struct {
	Result *float64     `json:"result"`
	Error  *errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestHandlerSuccess(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		requestBody string
		wantResult  float64
	}{
		{
			name:        "binary add",
			handler:     handler.Add,
			requestBody: `{"operands":[2,3]}`,
			wantResult:  5,
		},
		{
			name:        "binary subtract",
			handler:     handler.Subtract,
			requestBody: `{"operands":[10,4]}`,
			wantResult:  6,
		},
		{
			name:        "binary multiply",
			handler:     handler.Multiply,
			requestBody: `{"operands":[6,7]}`,
			wantResult:  42,
		},
		{
			name:        "binary divide",
			handler:     handler.Divide,
			requestBody: `{"operands":[20,4]}`,
			wantResult:  5,
		},
		{
			name:        "binary power",
			handler:     handler.Power,
			requestBody: `{"operands":[2,3]}`,
			wantResult:  8,
		},
		{
			name:        "unary square root",
			handler:     handler.SquareRoot,
			requestBody: `{"operands":[81]}`,
			wantResult:  9,
		},
		{
			name:        "unary percentage",
			handler:     handler.Percentage,
			requestBody: `{"operands":[25]}`,
			wantResult:  0.25,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := serve(t, test.handler, http.MethodPost, test.requestBody)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}

			body := decodeResponse(t, response)
			if body.Result == nil {
				t.Fatal("response is missing result")
			}
			if *body.Result != test.wantResult {
				t.Errorf("result = %v, want %v", *body.Result, test.wantResult)
			}
			if body.Error != nil {
				t.Errorf("unexpected error response: %+v", *body.Error)
			}
		})
	}
}

func TestHandlerRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		method      string
		requestBody string
		wantStatus  int
	}{
		{name: "add non-POST method", handler: handler.Add, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "subtract non-POST method", handler: handler.Subtract, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "multiply non-POST method", handler: handler.Multiply, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "divide non-POST method", handler: handler.Divide, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "power non-POST method", handler: handler.Power, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "square root non-POST method", handler: handler.SquareRoot, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "percentage non-POST method", handler: handler.Percentage, method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{
			name:        "malformed JSON",
			handler:     handler.Add,
			method:      http.MethodPost,
			requestBody: `{"operands":[2,3]`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "trailing JSON",
			handler:     handler.Add,
			method:      http.MethodPost,
			requestBody: `{"operands":[2,3]} {}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "malformed unary JSON",
			handler:     handler.Percentage,
			method:      http.MethodPost,
			requestBody: `{"operands":[25]`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "add invalid operand count",
			handler:     handler.Add,
			method:      http.MethodPost,
			requestBody: `{"operands":[2]}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "subtract invalid operand count",
			handler:     handler.Subtract,
			method:      http.MethodPost,
			requestBody: `{"operands":[2]}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "multiply invalid operand count",
			handler:     handler.Multiply,
			method:      http.MethodPost,
			requestBody: `{"operands":[2]}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "divide invalid operand count",
			handler:     handler.Divide,
			method:      http.MethodPost,
			requestBody: `{"operands":[2]}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "power invalid operand count",
			handler:     handler.Power,
			method:      http.MethodPost,
			requestBody: `{"operands":[2]}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "square root invalid operand count",
			handler:     handler.SquareRoot,
			method:      http.MethodPost,
			requestBody: `{"operands":[81,9]}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "percentage invalid operand count",
			handler:     handler.Percentage,
			method:      http.MethodPost,
			requestBody: `{"operands":[25,50]}`,
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := serve(t, test.handler, test.method, test.requestBody)
			assertErrorResponse(t, response, test.wantStatus, map[int]string{
				http.StatusBadRequest:       "INVALID_REQUEST",
				http.StatusMethodNotAllowed: "METHOD_NOT_ALLOWED",
			}[test.wantStatus])

			if test.wantStatus == http.StatusMethodNotAllowed {
				if allow := response.Header().Get("Allow"); allow != http.MethodPost {
					t.Errorf("Allow header = %q, want %q", allow, http.MethodPost)
				}
			}
		})
	}
}

func TestHandlerDomainErrors(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		requestBody string
		wantCode    string
	}{
		{
			name:        "division by zero",
			handler:     handler.Divide,
			requestBody: `{"operands":[20,0]}`,
			wantCode:    "DIVISION_BY_ZERO",
		},
		{
			name:        "negative square root",
			handler:     handler.SquareRoot,
			requestBody: `{"operands":[-1]}`,
			wantCode:    "NEGATIVE_SQUARE_ROOT",
		},
		{
			name:        "non-finite result",
			handler:     handler.Power,
			requestBody: `{"operands":[-1,0.5]}`,
			wantCode:    "INVALID_RESULT",
		},
		{
			name:        "non-finite add result",
			handler:     handler.Add,
			requestBody: `{"operands":[1e308,1e308]}`,
			wantCode:    "INVALID_RESULT",
		},
		{
			name:        "non-finite subtract result",
			handler:     handler.Subtract,
			requestBody: `{"operands":[-1e308,1e308]}`,
			wantCode:    "INVALID_RESULT",
		},
		{
			name:        "non-finite multiply result",
			handler:     handler.Multiply,
			requestBody: `{"operands":[1e308,2]}`,
			wantCode:    "INVALID_RESULT",
		},
		{
			name:        "non-finite divide result",
			handler:     handler.Divide,
			requestBody: `{"operands":[1e308,1e-308]}`,
			wantCode:    "INVALID_RESULT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := serve(t, test.handler, http.MethodPost, test.requestBody)
			assertErrorResponse(t, response, http.StatusUnprocessableEntity, test.wantCode)
		})
	}
}

func serve(t *testing.T, target http.HandlerFunc, method, body string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, "/", strings.NewReader(body))
	response := httptest.NewRecorder()
	target.ServeHTTP(response, request)
	return response
}

func decodeResponse(t *testing.T, response *httptest.ResponseRecorder) apiResponse {
	t.Helper()

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	var body apiResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body
}

func assertErrorResponse(
	t *testing.T,
	response *httptest.ResponseRecorder,
	wantStatus int,
	wantCode string,
) {
	t.Helper()

	if response.Code != wantStatus {
		t.Fatalf("status = %d, want %d", response.Code, wantStatus)
	}

	body := decodeResponse(t, response)
	if body.Error == nil {
		t.Fatal("response is missing error")
	}
	if body.Error.Code != wantCode {
		t.Errorf("error code = %q, want %q", body.Error.Code, wantCode)
	}
	if body.Error.Message == "" {
		t.Error("error message is empty")
	}
	if body.Result != nil {
		t.Errorf("unexpected result: %v", *body.Result)
	}
}
