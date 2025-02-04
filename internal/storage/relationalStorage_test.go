package storage

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestGetSwiftCodeDetails(t *testing.T) {
	t.Run("SwiftCodeNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		swiftCode := "INVALIDCODE"

		mock.ExpectQuery(`SELECT address, bankName, countryISO2, countryName, isHeadquarter, swiftCode FROM BanksData WHERE swiftCode = \$1`).
			WithArgs(swiftCode).
			WillReturnError(sql.ErrNoRows)

		result, err := storage.GetSwiftCodeDetails(swiftCode)
		if !errors.Is(err, ErrSwiftCodeNotFound) {
			t.Errorf("expected error %v, got %v", ErrSwiftCodeNotFound, err)
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("NonHeadquarterBank", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		swiftCode := "TESTPL33AAA"

		mock.ExpectQuery(`SELECT address, bankName, countryISO2, countryName, isHeadquarter, swiftCode FROM BanksData WHERE swiftCode = \$1`).
			WithArgs(swiftCode).
			WillReturnRows(sqlmock.NewRows([]string{"address", "bankName", "countryISO2", "countryName", "isHeadquarter", "swiftCode"}).
				AddRow("Address", "Bank", "PL", "POLAND", false, swiftCode))

		result, err := storage.GetSwiftCodeDetails(swiftCode)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected bank details")
		}
		if len(result.Branches) > 0 {
			t.Error("expected no branches for non-headquarter bank")
		}
	})

	t.Run("HeadquarterWithBranches", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		swiftCode := "TESTPL33XXX"

		mock.ExpectQuery(`SELECT address, bankName, countryISO2, countryName, isHeadquarter, swiftCode FROM BanksData WHERE swiftCode = \$1`).
			WithArgs(swiftCode).
			WillReturnRows(sqlmock.NewRows([]string{"address", "bankName", "countryISO2", "countryName", "isHeadquarter", "swiftCode"}).
				AddRow("HQ Address", "HQ Bank", "PL", "POLAND", true, swiftCode))

		likeParam := swiftCode[:8] + "%"
		mock.ExpectQuery(`SELECT address, bankName, countryISO2, isHeadquarter, swiftCode FROM BanksData WHERE swiftCode LIKE \$1`).
			WithArgs(likeParam).
			WillReturnRows(sqlmock.NewRows([]string{"address", "bankName", "countryISO2", "isHeadquarter", "swiftCode"}).
				AddRow("Branch Address", "Branch Bank", "PL", false, "HEADQCODE123").
				AddRow("Other Branch", "Other Bank", "PL", false, "HEADQCODE456"))

		result, err := storage.GetSwiftCodeDetails(swiftCode)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected bank details")
		}
		if len(result.Branches) != 2 {
			t.Errorf("expected 2 branches, got %d", len(result.Branches))
		}
	})
}

func TestGetSwiftCodesForCountry(t *testing.T) {
	t.Run("ValidCountryWithoutData", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		iso2Code := "PL"

		mock.ExpectQuery(`SELECT countryISO2, countryName, address, bankName, isHeadquarter, swiftCode FROM BanksData WHERE countryISO2 = \$1`).
			WithArgs(iso2Code).
			WillReturnRows(sqlmock.NewRows([]string{}))

		result, err := storage.GetSwiftCodesForCountry(iso2Code)
		if !errors.Is(err, ErrISO2CodeNotFound) {
			t.Errorf("expected error %v, got %v", ErrISO2CodeNotFound, err)
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("ValidCountry", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		iso2Code := "US"

		mock.ExpectQuery(`SELECT countryISO2, countryName, address, bankName, isHeadquarter, swiftCode FROM BanksData WHERE countryISO2 = \$1`).
			WithArgs(iso2Code).
			WillReturnRows(sqlmock.NewRows([]string{"countryISO2", "countryName", "address", "bankName", "isHeadquarter", "swiftCode"}).
				AddRow("US", "USA", "Addr1", "Bank1", false, "BANKUS11XXX").
				AddRow("US", "USA", "Addr2", "Bank2", true, "BANKUS22XXX"))

		result, err := storage.GetSwiftCodesForCountry(iso2Code)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected country banks details")
		}
		if len(result.SwiftCodes) != 2 {
			t.Errorf("expected 2 swift codes, got %d", len(result.SwiftCodes))
		}
	})
}

func TestAddSwiftCodeEntry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		bank := Bank{
			Address:       strPtr("New Address"),
			BankName:      strPtr("New Bank"),
			CountryISO2:   strPtr("PL"),
			CountryName:   strPtr("POLAND"),
			IsHeadquarter: boolPtr(true),
			SwiftCode:     strPtr("TESTPL33XXX"),
		}

		mock.ExpectExec(`INSERT INTO BanksData \(address, bankName, countryISO2, countryName, isHeadquarter, swiftCode\) 
VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
			WithArgs(bank.Address, bank.BankName, bank.CountryISO2, bank.CountryName, bank.IsHeadquarter, bank.SwiftCode).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = storage.AddSwiftCodeEntry(bank)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("DuplicateEntry", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		bank := Bank{SwiftCode: strPtr("DUPLICATE")}

		mock.ExpectExec(`INSERT INTO BanksData .+`).
			WillReturnError(errors.New("duplicate key value violates unique constraint"))

		err = storage.AddSwiftCodeEntry(bank)
		if !errors.Is(err, ErrSwiftCodeExists) {
			t.Errorf("expected error %v, got %v", ErrSwiftCodeExists, err)
		}
	})
}

func TestDeleteSwiftCodeEntry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		swiftCode := "TODELETE"

		mock.ExpectExec(`SELECT swiftCode FROM BanksData WHERE swiftCode = \$1`).
			WithArgs(swiftCode).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(`DELETE FROM BanksData WHERE swiftCode = \$1`).
			WithArgs(swiftCode).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = storage.DeleteSwiftCodeEntry(swiftCode)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock: %v", err)
		}
		defer db.Close()

		storage := NewRelationalDB(db)
		swiftCode := "NOTFOUND"

		mock.ExpectExec(`SELECT swiftCode FROM BanksData WHERE swiftCode = \$1`).
			WithArgs(swiftCode).
			WillReturnError(sql.ErrNoRows)

		err = storage.DeleteSwiftCodeEntry(swiftCode)
		if !errors.Is(err, ErrSwiftCodeNotFound) {
			t.Errorf("expected error %v, got %v", ErrSwiftCodeNotFound, err)
		}
	})
}
