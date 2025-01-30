package app

import (
	"encoding/json"
	"errors"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock Storage implementing the storage.Storage interface
type mockStorage struct {
	GetSwiftCodeDetailsFunc     func(swiftCode string) (*storage.Bank, error)
	GetSwiftCodesForCountryFunc func(iso2Code string) (*storage.CountryBanks, error)
	AddSwiftCodeEntryFunc       func(b storage.Bank) error
	DeleteSwiftCodeEntryFunc    func(swiftCode string) error
}

func (m *mockStorage) GetSwiftCodeDetails(swiftCode string) (*storage.Bank, error) {
	return m.GetSwiftCodeDetailsFunc(swiftCode)
}

func (m *mockStorage) GetSwiftCodesForCountry(iso2Code string) (*storage.CountryBanks, error) {
	return m.GetSwiftCodesForCountryFunc(iso2Code)
}

func (m *mockStorage) AddSwiftCodeEntry(b storage.Bank) error {
	return m.AddSwiftCodeEntryFunc(b)
}

func (m *mockStorage) DeleteSwiftCodeEntry(swiftCode string) error {
	return m.DeleteSwiftCodeEntryFunc(swiftCode)
}

// Helper to set path variables in the request context
func setPathVars(r *http.Request, vars map[string]string) *http.Request {
	for k, v := range vars {
		r.SetPathValue(k, v)
	}
	return r
}

func TestHandleGetSwiftCodeDetails(t *testing.T) {
	tests := []struct {
		name           string
		swiftCode      string
		mockStorage    *mockStorage
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "empty swift code",
			swiftCode:      "",
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "swiftCode not found in path",
		},
		{
			name:      "swift code not found",
			swiftCode: "INVALID_CODE",
			mockStorage: &mockStorage{
				GetSwiftCodeDetailsFunc: func(_ string) (*storage.Bank, error) {
					return nil, storage.ErrSwiftCodeNotFound
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "swiftCode is invalid",
		},
		{
			name:      "swift code not found",
			swiftCode: "TESTPL33XXX",
			mockStorage: &mockStorage{
				GetSwiftCodeDetailsFunc: func(_ string) (*storage.Bank, error) {
					return nil, storage.ErrSwiftCodeNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    storage.ErrSwiftCodeNotFound.Error(),
		},
		{
			name:      "storage error",
			swiftCode: "TESTPL33XXX",
			mockStorage: &mockStorage{
				GetSwiftCodeDetailsFunc: func(_ string) (*storage.Bank, error) {
					return nil, errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "storage error",
		},
		{
			name:      "success",
			swiftCode: "TESTPL33XXX",
			mockStorage: &mockStorage{
				GetSwiftCodeDetailsFunc: func(swiftCode string) (*storage.Bank, error) {
					return &storage.Bank{SwiftCode: &swiftCode}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/swift-codes/"+tt.swiftCode, nil)
			req = setPathVars(req, map[string]string{"swiftCode": tt.swiftCode})

			rec := httptest.NewRecorder()

			service := NewBankService(tt.mockStorage)
			service.handleGetSwiftCodeDetails(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var bank storage.Bank
				assert.NoError(t, json.NewDecoder(res.Body).Decode(&bank))
				assert.Equal(t, tt.swiftCode, *bank.SwiftCode)
			} else {
				var resp storage.Response
				assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
				assert.Equal(t, tt.expectedMsg, resp.Message)
			}
		})
	}
}

func TestHandleGetCountrySwiftCodes(t *testing.T) {
	tests := []struct {
		name           string
		countryCode    string
		mockStorage    *mockStorage
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "empty country code",
			countryCode:    "",
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "countryISO2code not found in path",
		},
		{
			name:           "lowercase country code",
			countryCode:    "us",
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "countryISO2code is invalid",
		},
		{
			name:           "invalid length",
			countryCode:    "USA",
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "countryISO2code is invalid",
		},
		{
			name:        "valid code not found",
			countryCode: "XX",
			mockStorage: &mockStorage{
				GetSwiftCodesForCountryFunc: func(_ string) (*storage.CountryBanks, error) {
					return nil, storage.ErrISO2CodeNotFound
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "countryISO2code is invalid",
		},
		{
			name:        "storage error",
			countryCode: "US",
			mockStorage: &mockStorage{
				GetSwiftCodesForCountryFunc: func(_ string) (*storage.CountryBanks, error) {
					return nil, errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "storage error",
		},
		{
			name:        "iso2 code not in database",
			countryCode: "US",
			mockStorage: &mockStorage{
				GetSwiftCodesForCountryFunc: func(_ string) (*storage.CountryBanks, error) {
					return nil, storage.ErrISO2CodeNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    storage.ErrISO2CodeNotFound.Error(),
		},
		{
			name:        "success",
			countryCode: "US",
			mockStorage: &mockStorage{
				GetSwiftCodesForCountryFunc: func(iso string) (*storage.CountryBanks, error) {
					return &storage.CountryBanks{CountryISO2: iso}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/swift-codes/country/"+tt.countryCode, nil)
			req = setPathVars(req, map[string]string{"countryISO2code": tt.countryCode})

			rec := httptest.NewRecorder()

			service := NewBankService(tt.mockStorage)
			service.handleGetCountrySwiftCodes(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var countryBanks storage.CountryBanks
				assert.NoError(t, json.NewDecoder(res.Body).Decode(&countryBanks))
				assert.Equal(t, tt.countryCode, countryBanks.CountryISO2)
			} else {
				var resp storage.Response
				assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
				assert.Equal(t, tt.expectedMsg, resp.Message)
			}
		})
	}
}

func TestHandleAddSwiftCodeDetails(t *testing.T) {
	validBank := `{
		"address": "123 Test St",
		"bankName": "Test Bank",
		"countryISO2": "PL",
		"countryName": "POLAND",
		"isHeadquarter": true,
		"swiftCode": "TESTPL33XXX",
		"branches": []
	}`

	tests := []struct {
		name           string
		body           io.Reader
		mockStorage    *mockStorage
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "invalid JSON",
			body:           strings.NewReader(`{invalid}`),
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Error parsing request body",
		},
		{
			name: "validation error (missing fields)",
			body: strings.NewReader(`{"bankName": "No Swift Code"}`),
			mockStorage: &mockStorage{
				AddSwiftCodeEntryFunc: func(_ storage.Bank) error {
					return nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "address is required",
		},
		{
			name: "duplicate swift code",
			body: strings.NewReader(validBank),
			mockStorage: &mockStorage{
				AddSwiftCodeEntryFunc: func(_ storage.Bank) error {
					return storage.ErrSwiftCodeExists
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    storage.ErrSwiftCodeExists.Error(),
		},
		{
			name: "storage error",
			body: strings.NewReader(validBank),
			mockStorage: &mockStorage{
				AddSwiftCodeEntryFunc: func(_ storage.Bank) error {
					return errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "storage error",
		},
		{
			name: "success",
			body: strings.NewReader(validBank),
			mockStorage: &mockStorage{
				AddSwiftCodeEntryFunc: func(_ storage.Bank) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Successfully added bank with swift code TESTPL33XXX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/swift-codes", tt.body)
			rec := httptest.NewRecorder()

			service := NewBankService(tt.mockStorage)
			service.handleAddSwiftCodeDetails(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			var resp storage.Response
			assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
			assert.Equal(t, tt.expectedMsg, resp.Message)
		})
	}
}

func TestHandleDeleteSwiftCode(t *testing.T) {
	tests := []struct {
		name           string
		swiftCode      string
		mockStorage    *mockStorage
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "empty swift code",
			swiftCode:      "",
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "swift-code not found in path",
		},
		{
			name:           "invalid swift code",
			swiftCode:      "invalid",
			mockStorage:    &mockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "swift-code is invalid",
		},
		{
			name:      "swift code not found",
			swiftCode: "TESTPL33XXX",
			mockStorage: &mockStorage{
				DeleteSwiftCodeEntryFunc: func(_ string) error {
					return storage.ErrSwiftCodeNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    storage.ErrSwiftCodeNotFound.Error(),
		},
		{
			name:      "storage error",
			swiftCode: "TESTPL33XXX",
			mockStorage: &mockStorage{
				DeleteSwiftCodeEntryFunc: func(_ string) error {
					return errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "storage error",
		},
		{
			name:      "success",
			swiftCode: "TESTPL33XXX",
			mockStorage: &mockStorage{
				DeleteSwiftCodeEntryFunc: func(_ string) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Bank with swift code: TESTPL33XXX has been deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/swift-codes/"+tt.swiftCode, nil)
			req = setPathVars(req, map[string]string{"swiftCode": tt.swiftCode})

			rec := httptest.NewRecorder()

			service := NewBankService(tt.mockStorage)
			service.handleDeleteSwiftCode(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			var resp storage.Response
			assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
			assert.Equal(t, tt.expectedMsg, resp.Message)
		})
	}
}
