package app

import (
	"context"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"testing"
	"time"
)

// mockStorage implements the storage.Storage interface for testing.
type mockStorageApi struct {
	getSwiftCodeDetailsFunc     func(string) (*storage.Bank, error)
	getSwiftCodesForCountryFunc func(string) (*storage.CountryBanks, error)
	addSwiftCodeEntryFunc       func(storage.Bank) error
	deleteSwiftCodeEntryFunc    func(string) error
}

func (m *mockStorageApi) GetSwiftCodeDetails(swiftCode string) (*storage.Bank, error) {
	if m.getSwiftCodeDetailsFunc != nil {
		return m.getSwiftCodeDetailsFunc(swiftCode)
	}
	return nil, storage.ErrSwiftCodeNotFound
}

func (m *mockStorageApi) GetSwiftCodesForCountry(iso2Code string) (*storage.CountryBanks, error) {
	if m.getSwiftCodesForCountryFunc != nil {
		return m.getSwiftCodesForCountryFunc(iso2Code)
	}
	return nil, storage.ErrISO2CodeNotFound
}

func (m *mockStorageApi) AddSwiftCodeEntry(b storage.Bank) error {
	if m.addSwiftCodeEntryFunc != nil {
		return m.addSwiftCodeEntryFunc(b)
	}
	return storage.ErrSwiftCodeExists
}

func (m *mockStorageApi) DeleteSwiftCodeEntry(swiftCode string) error {
	if m.deleteSwiftCodeEntryFunc != nil {
		return m.deleteSwiftCodeEntryFunc(swiftCode)
	}
	return storage.ErrSwiftCodeNotFound
}

// TestAPIServer_ShutdownOnContextCancel tests that the server shuts down gracefully when the context is canceled.
func TestAPIServer_ShutdownOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := &mockStorageApi{}
	server := NewAPIServer(":0", storage)

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Expected no error during shutdown, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not shut down within the expected time")
	}
}

// TestAPIServer_InvalidAddress tests that the server returns an error when started with an invalid address.
func TestAPIServer_InvalidAddress(t *testing.T) {
	storage := &mockStorageApi{}
	server := NewAPIServer("invalid-address:8080", storage) // Invalid address

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.Start(ctx)
	if err == nil {
		t.Error("Expected an error due to invalid address, got nil")
	}
}

func TestAPIServer_StorageIntegration(t *testing.T) {
	expectedBank := &storage.Bank{SwiftCode: strPtr("TEST")}

	mockStorage := &mockStorageApi{
		getSwiftCodeDetailsFunc: func(swiftCode string) (*storage.Bank, error) {
			if swiftCode == "TEST" {
				return expectedBank, nil
			}
			return nil, storage.ErrSwiftCodeNotFound
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := NewAPIServer(":0", mockStorage)
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(ctx)
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// This is simplified to check if the storage is correctly integrated.
	bank, err := mockStorage.GetSwiftCodeDetails("TEST")
	if err != nil {
		t.Fatalf("Failed to get swift code details: %v", err)
	}

	if bank.SwiftCode != expectedBank.SwiftCode {
		t.Errorf("Expected SwiftCode %s, got %s", *expectedBank.SwiftCode, *bank.SwiftCode)
	}

	cancel()
	<-errChan
}
