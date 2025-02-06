//go:build unit

package utils

import (
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}

	strPtr := func(s string) *string { return &s }
	boolPtr := func(b bool) *bool { return &b }

	tests := []struct {
		name           string
		status         int
		data           any
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Write simple JSON response",
			status:         http.StatusOK,
			data:           response{Message: "Success"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Success"}`,
		},
		{
			name:           "Write error JSON response",
			status:         http.StatusBadRequest,
			data:           response{Message: "Bad Request"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Bad Request"}`,
		},
		{
			name:   "Write BankData as JSON response",
			status: http.StatusOK,
			data: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"address":"123 Test Street","bankName":"Test Bank","countryISO2":"PL","countryName":"POLAND","isHeadquarter":true,"swiftCode":"TESTPL33XXX","branches":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			WriteJSON(recorder, tt.status, tt.data)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" && tt.data != nil {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			body := strings.TrimSpace(recorder.Body.String())
			if body != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}
