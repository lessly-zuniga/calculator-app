package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type calculatorResponse struct {
	Result *float64      `json:"result"`
	Error  *errorPayload `json:"error"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestCalculatorEndpointsSuccess(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		want float64
	}{
		{name: "add", path: "/api/v1/add", body: `{"operands":[2,3]}`, want: 5},
		{name: "subtract", path: "/api/v1/subtract", body: `{"operands":[10,4]}`, want: 6},
		{name: "multiply", path: "/api/v1/multiply", body: `{"operands":[6,7]}`, want: 42},
		{name: "divide", path: "/api/v1/divide", body: `{"operands":[20,4]}`, want: 5},
		{name: "square root", path: "/api/v1/square-root", body: `{"operands":[81]}`, want: 9},
		{name: "percentage", path: "/api/v1/percentage", body: `{"operands":[25]}`, want: 0.25},
		{name: "power", path: "/api/v1/power", body: `{"operands":[2,3]}`, want: 8},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := performRequest(t, http.MethodPost, test.path, test.body)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			assertJSONContentType(t, response)

			var body calculatorResponse
			decodeResponse(t, response, &body)
			if body.Result == nil {
				t.Fatal("response is missing result")
			}
			if *body.Result != test.want {
				t.Errorf("result = %v, want %v", *body.Result, test.want)
			}
			if body.Error != nil {
				t.Errorf("unexpected error response: %+v", *body.Error)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	t.Run("GET succeeds", func(t *testing.T) {
		response := performRequest(t, http.MethodGet, "/api/v1/health", "")
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
		}
		assertJSONContentType(t, response)

		var body struct {
			Status string `json:"status"`
		}
		decodeResponse(t, response, &body)
		if body.Status != "ok" {
			t.Errorf("status value = %q, want %q", body.Status, "ok")
		}
	})

	t.Run("non-GET is rejected", func(t *testing.T) {
		response := performRequest(t, http.MethodPost, "/api/v1/health", "")
		assertErrorResponse(t, response, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET requests are allowed")
		if allow := response.Header().Get("Allow"); allow != http.MethodGet {
			t.Errorf("Allow header = %q, want %q", allow, http.MethodGet)
		}
	})
}

func TestCalculatorEndpointValidation(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		body    string
		message string
	}{
		{name: "malformed JSON", path: "/api/v1/add", body: `{"operands":[2,3]`, message: "Exactly two numeric operands are required"},
		{name: "missing operands", path: "/api/v1/subtract", body: `{}`, message: "Exactly two numeric operands are required"},
		{name: "incorrect binary operand count", path: "/api/v1/multiply", body: `{"operands":[6]}`, message: "Exactly two numeric operands are required"},
		{name: "incorrect unary operand count", path: "/api/v1/percentage", body: `{"operands":[25,50]}`, message: "Exactly one numeric operand is required"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := performRequest(t, http.MethodPost, test.path, test.body)
			assertErrorResponse(t, response, http.StatusBadRequest, "INVALID_REQUEST", test.message)
		})
	}
}

func TestCalculatorEndpointsRejectNonPOSTMethods(t *testing.T) {
	paths := []string{
		"/api/v1/add",
		"/api/v1/subtract",
		"/api/v1/multiply",
		"/api/v1/divide",
		"/api/v1/square-root",
		"/api/v1/percentage",
		"/api/v1/power",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			response := performRequest(t, http.MethodGet, path, "")
			assertErrorResponse(t, response, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST requests are allowed")
			if allow := response.Header().Get("Allow"); allow != http.MethodPost {
				t.Errorf("Allow header = %q, want %q", allow, http.MethodPost)
			}
		})
	}
}

func TestCalculatorEndpointDomainErrors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		body    string
		code    string
		message string
	}{
		{
			name:    "division by zero",
			path:    "/api/v1/divide",
			body:    `{"operands":[20,0]}`,
			code:    "DIVISION_BY_ZERO",
			message: "Cannot divide by zero",
		},
		{
			name:    "negative square root",
			path:    "/api/v1/square-root",
			body:    `{"operands":[-1]}`,
			code:    "NEGATIVE_SQUARE_ROOT",
			message: "Cannot calculate the square root of a negative number",
		},
		{
			name:    "non-finite power result",
			path:    "/api/v1/power",
			body:    `{"operands":[-1,0.5]}`,
			code:    "INVALID_RESULT",
			message: "Calculation produced an invalid numeric result",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := performRequest(t, http.MethodPost, test.path, test.body)
			assertErrorResponse(t, response, http.StatusUnprocessableEntity, test.code, test.message)
		})
	}
}

func performRequest(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	newMux().ServeHTTP(response, request)
	return response
}

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantCode, wantMessage string) {
	t.Helper()

	if response.Code != wantStatus {
		t.Fatalf("status = %d, want %d", response.Code, wantStatus)
	}
	assertJSONContentType(t, response)

	var body calculatorResponse
	decodeResponse(t, response, &body)
	if body.Error == nil {
		t.Fatal("response is missing error")
	}
	if body.Error.Code != wantCode {
		t.Errorf("error code = %q, want %q", body.Error.Code, wantCode)
	}
	if body.Error.Message != wantMessage {
		t.Errorf("error message = %q, want %q", body.Error.Message, wantMessage)
	}
	if body.Result != nil {
		t.Errorf("unexpected result: %v", *body.Result)
	}
}

func assertJSONContentType(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}

func decodeResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()

	if err := json.NewDecoder(response.Body).Decode(destination); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
