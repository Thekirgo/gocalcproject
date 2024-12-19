package tests

import (
	"bytes"
	"calculator-service/internal/api"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculatorHandler_Calculate(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		wantStatus     int
		wantResult     *float64
		wantErrMessage string
	}{
		{
			name: "valid expression",
			requestBody: api.CalculateRequest{
				Expression: "2+2*2",
			},
			wantStatus: http.StatusOK,
			wantResult: func() *float64 { f := 6.0; return &f }(),
		},
		{
			name: "invalid expression",
			requestBody: api.CalculateRequest{
				Expression: "2+a",
			},
			wantStatus:     http.StatusUnprocessableEntity,
			wantErrMessage: "Expression is not valid",
		},
		{
			name:           "empty request",
			requestBody:    nil,
			wantStatus:     http.StatusInternalServerError,
			wantErrMessage: "Internal server error",
		},
		{
			name: "empty expression",
			requestBody: api.CalculateRequest{
				Expression: "",
			},
			wantStatus:     http.StatusUnprocessableEntity,
			wantErrMessage: "Expression is not valid",
		},
		{
			name: "division by zero",
			requestBody: api.CalculateRequest{
				Expression: "1/0",
			},
			wantStatus:     http.StatusUnprocessableEntity,
			wantErrMessage: "Expression is not valid",
		},
	}

	handler := api.NewCalculatorHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.Calculate(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Calculate() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantResult != nil {
				var response api.SuccessResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				if response.Result != *tt.wantResult {
					t.Errorf("Calculate() result = %v, want %v", response.Result, *tt.wantResult)
				}
			}

			if tt.wantErrMessage != "" {
				var response api.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				if response.Error != tt.wantErrMessage {
					t.Errorf("Calculate() error = %v, want %v", response.Error, tt.wantErrMessage)
				}
			}
		})
	}
}