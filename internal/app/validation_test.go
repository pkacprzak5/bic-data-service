//go:build unit

package app

import (
	"errors"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestValidateBankData_OneMissingValue(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name: "Missing address",
			input: storage.Bank{
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing bank name",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing countryISO2",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing country name",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryName is required"),
		},
		{
			name: "Missing isHeadquarter",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("isHeadquarter is required"),
		},
		{
			name: "Missing swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("swiftCode is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankData_TwoMissingValues(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name: "Missing address and bank name",
			input: storage.Bank{
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address and countryISO2",
			input: storage.Bank{
				BankName:      strPtr("Test Bank"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address and country name",
			input: storage.Bank{
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address and isHeadquarter",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address and swiftCode",
			input: storage.Bank{
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing bank name and countryISO2",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name and country name",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name and isHeadquarter",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing countryISO2 and country name",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing countryISO2 and isHeadquarter",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				BankName:    strPtr("Test Bank"),
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing countryISO2 and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing country name and isHeadquarter",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryName is required"),
		},
		{
			name: "Missing country name and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("countryName is required"),
		},
		{
			name: "Missing isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("isHeadquarter is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankData_ThreeMissingValues(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name: "Missing address, bank name, and countryISO2",
			input: storage.Bank{
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, and country name",
			input: storage.Bank{
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, and isHeadquarter",
			input: storage.Bank{
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, and swiftCode",
			input: storage.Bank{
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2 and country name",
			input: storage.Bank{
				BankName:      strPtr("Test Bank"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2 and isHeadquarter",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2 and swiftCode",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, country name and isHeadquarter",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, country name and swiftCode",
			input: storage.Bank{
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, isHeadquarter and swiftCode",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing bank name, countryISO2 and country name",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, countryISO2 and isHeadquarter",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, countryISO2 and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, country name and isHeadquarter",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				CountryISO2: strPtr("PL"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, country name and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing countryISO2, country name and isHeadquarter",
			input: storage.Bank{
				Address:   strPtr("123 Test Street"),
				BankName:  strPtr("Test Bank"),
				SwiftCode: strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing countryISO2, country name and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing countryISO2, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				BankName:    strPtr("Test Bank"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("countryISO2 is required"),
		},
		{
			name: "Missing country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
			},
			wantError: errors.New("countryName is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankData_FourMissingValues(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name: "Missing address, bank name, countryISO2 and country name",
			input: storage.Bank{
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, countryISO2 and isHeadquarter",
			input: storage.Bank{
				CountryName: strPtr("POLAND"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, countryISO2 and swiftCode",
			input: storage.Bank{
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, country name and isHeadquarter",
			input: storage.Bank{
				CountryISO2: strPtr("PL"),
				SwiftCode:   strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, country name and swiftCode",
			input: storage.Bank{
				CountryISO2:   strPtr("PL"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, isHeadquarter and swiftCode",
			input: storage.Bank{
				CountryISO2: strPtr("PL"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2, country name and isHeadquarter",
			input: storage.Bank{
				BankName:  strPtr("Test Bank"),
				SwiftCode: strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2, country name and swiftCode",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2, isHeadquarter and swiftCode",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				BankName:    strPtr("Test Bank"),
				CountryISO2: strPtr("PL"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing bank name, countryISO2, country name and isHeadquarter",
			input: storage.Bank{
				Address:   strPtr("123 Test Street"),
				SwiftCode: strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, countryISO2, country name and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, countryISO2, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing bank name, country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:     strPtr("123 Test Street"),
				CountryISO2: strPtr("PL"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Missing countryISO2, country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:  strPtr("123 Test Street"),
				BankName: strPtr("Test Bank"),
			},
			wantError: errors.New("countryISO2 is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankName_FiveMissingValues(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name: "Missing address, bank name, countryISO2, country name and isHeadquarter",
			input: storage.Bank{
				SwiftCode: strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, countryISO2, country name and swiftCode",
			input: storage.Bank{
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, countryISO2, isHeadquarter and swiftCode",
			input: storage.Bank{
				CountryName: strPtr("POLAND"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, bank name, country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				CountryISO2: strPtr("PL"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing address, countryISO2, country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				BankName: strPtr("Test Bank"),
			},
			wantError: errors.New("address is required"),
		},
		{
			name: "Missing bank name, countryISO2, country name, isHeadquarter and swiftCode",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				IsHeadquarter: boolPtr(true),
			},
			wantError: errors.New("bankName is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankData_AllValuesMissing(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name:      "All fields missing",
			input:     storage.Bank{},
			wantError: errors.New("address is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankData_InvalidData(t *testing.T) {
	tests := []struct {
		name      string
		input     storage.Bank
		wantError error
	}{
		{
			name: "Bank name is empty",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr(""),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Bank name have only whitespaces",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("  "),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("bankName is required"),
		},
		{
			name: "Invalid ISO2 code",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("ZZ"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTZZ33XXX"),
			},
			wantError: errors.New("countryISO2 is invalid"),
		},
		{
			name: "Country name does not match ISO2 code",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("GERMANY"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("countryName does not match ISO2 code"),
		},
		{
			name: "Swift code does not match ISO2 code",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTZZ33XXX"),
			},
			wantError: errors.New("swiftCode is invalid"),
		},
		{
			name: "Invalid SWIFT code format",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("INVALIDSWIFT"),
			},
			wantError: errors.New("swiftCode is invalid"),
		},
		{
			name: "SWIFT code indicates branch but marked as headquarter",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33AAA"),
			},
			wantError: errors.New("swiftCode indicates bank's branch"),
		},
		{
			name: "SWIFT code indicates headquarter but marked as branch",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(false),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
			wantError: errors.New("swiftCode indicates bank's headquarter"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err == nil || err.Error() != tt.wantError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateBankData_CorrectData(t *testing.T) {
	tests := []struct {
		name  string
		input storage.Bank
	}{
		{
			name: "Valid bank data - headquarters",
			input: storage.Bank{
				Address:       strPtr("123 Test Street"),
				BankName:      strPtr("Test Bank"),
				CountryISO2:   strPtr("PL"),
				CountryName:   strPtr("POLAND"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("TESTPL33XXX"),
			},
		},
		{
			name: "Valid bank data - branch",
			input: storage.Bank{
				Address:       strPtr("456 Another Street"),
				BankName:      strPtr("Another Bank"),
				CountryISO2:   strPtr("DE"),
				CountryName:   strPtr("GERMANY"),
				IsHeadquarter: boolPtr(false),
				SwiftCode:     strPtr("ANOTDE33AAA"),
			},
		},
		{
			name: "Valid bank data - US headquarters",
			input: storage.Bank{
				Address:       strPtr("1600 Pennsylvania Avenue NW"),
				BankName:      strPtr("USA National Bank"),
				CountryISO2:   strPtr("US"),
				CountryName:   strPtr("UNITED STATES OF AMERICA (THE)"),
				IsHeadquarter: boolPtr(true),
				SwiftCode:     strPtr("USNBUS33XXX"),
			},
		},
		{
			name: "Valid bank data - Italian branch",
			input: storage.Bank{
				Address:       strPtr("Via Roma 123"),
				BankName:      strPtr("Italian Bank"),
				CountryISO2:   strPtr("IT"),
				CountryName:   strPtr("ITALY"),
				IsHeadquarter: boolPtr(false),
				SwiftCode:     strPtr("ITBKIT44BBB"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBankData(tt.input)
			if err != nil {
				t.Errorf("unexpected error for valid data: %v", err)
			}
		})
	}
}

func TestIsValidISO2_ValidCodes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Valid ISO2Code for Poland",
			input: "PL",
			want:  true,
		},
		{
			name:  "Valid ISO2Code for USA",
			input: "US",
			want:  true,
		},
		{
			name:  "Valid ISO2Code for Germany",
			input: "DE",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := isValidISO2(tt.input)
			if !flag {
				t.Errorf("isValidISO2(%v) returned false, should return true", tt.input)
			}
		})
	}
}

func TestIsValidISO2_InvalidCodes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Invalid ISO2Code ZZ",
			input: "ZZ",
			want:  false,
		},
		{
			name:  "Invalid ISO2Code pl",
			input: "pl",
			want:  false,
		},
		{
			name:  "Invalid ISO2Code pl",
			input: "1234",
			want:  false,
		},
		{
			name:  "Invalid ISO2Code pl",
			input: "PLDE",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := isValidISO2(tt.input)
			if flag {
				t.Errorf("isValidISO2(%v) returned true, should return false", tt.input)
			}
		})
	}
}

func TestIso2CodeToCountry(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ISO2Code for Poland",
			input: "PL",
			want:  "POLAND",
		},
		{
			name:  "ISO2Code for Germany",
			input: "DE",
			want:  "GERMANY",
		},
		{
			name:  "Invalid ISO2Code ZZ",
			input: "ZZ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countryName := iso2CodeToCountry(tt.input)
			if countryName != tt.want {
				t.Errorf("iso2CodeToCountry(%v) returned %v, should return %v", tt.input, countryName, tt.want)
			}
		})
	}
}

func TestIsValidSWIFT(t *testing.T) {
	tests := []struct {
		name      string
		inputCode string
		iso2Code  string
		want      bool
	}{
		{
			name:      "Valid SWIFT Code for Poland",
			inputCode: "TESTPL33AAA",
			iso2Code:  "PL",
			want:      true,
		},
		{
			name:      "Valid SWIFT Code for Germany",
			inputCode: "TESTDE33XXX",
			iso2Code:  "DE",
			want:      true,
		},
		{
			name:      "Invalid SWIFT Code for Poland",
			inputCode: "TESTDE33XXX",
			iso2Code:  "PL",
			want:      false,
		},
		{
			name:      "Invalid SWIFT Code for Germany",
			inputCode: "TESTPL33XXX",
			iso2Code:  "DE",
			want:      false,
		},
		{
			name:      "Invalid SWIFT Code format (only numbers)",
			inputCode: "1234",
			iso2Code:  "PL",
			want:      false,
		},
		{
			name:      "Invalid SWIFT Code format (wrong length)",
			inputCode: "AABBDDSSGGGGSS",
			iso2Code:  "PL",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := isValidSWIFT(tt.inputCode, tt.iso2Code)
			if flag != tt.want {
				t.Errorf("isValidSWIFT(%v) returned %v, should return %v", tt.inputCode, flag, tt.want)
			}
		})
	}
}
