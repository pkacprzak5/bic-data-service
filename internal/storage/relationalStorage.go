package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type RelationalDB struct {
	db *sql.DB
}

func NewRelationalDB(db *sql.DB) *RelationalDB {
	return &RelationalDB{db: db}
}

func (r *RelationalDB) GetSwiftCodeDetails(swiftCode string) (*Bank, error) {
	query := `SELECT address, bankName, countryISO2, countryName, isHeadquarter, swiftCode
		FROM BanksData
		WHERE swiftCode = $1`

	var bank Bank
	err := r.db.QueryRow(query, swiftCode).
		Scan(&bank.Address, &bank.BankName, &bank.CountryISO2, &bank.CountryName, &bank.IsHeadquarter, &bank.SwiftCode)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSwiftCodeNotFound
	} else if err != nil {
		return nil, err
	}

	if !*bank.IsHeadquarter {
		return &bank, nil
	}

	query = `SELECT address, bankName, countryISO2, isHeadquarter, swiftCode
		FROM BanksData
		WHERE swiftCode LIKE $1`

	rows, err := r.db.Query(query, fmt.Sprintf("%s%%", swiftCode[:8]))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b BankBranch

		err := rows.Scan(&b.Address, &b.BankName, &b.CountryISO2, &b.IsHeadquarter, &b.SwiftCode)
		if err != nil {
			return nil, err
		}

		if b.IsHeadquarter {
			continue
		}
		bank.Branches = append(bank.Branches, b)
	}

	return &bank, nil
}

func (r *RelationalDB) GetSwiftCodesForCountry(iso2Code string) (*CountryBanks, error) {
	query := `SELECT countryISO2, countryName, address, bankName, isHeadquarter, swiftCode 
		FROM BanksData
		WHERE countryISO2 = $1`

	rows, err := r.db.Query(query, iso2Code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrISO2CodeNotFound
	}

	var countryBanks CountryBanks
	for {
		var address, bankName, countryName, swiftCode, countryISO2 string
		var isHeadquarter bool

		err := rows.Scan(&countryISO2, &countryName, &address, &bankName, &isHeadquarter, &swiftCode)
		if err != nil {
			return nil, err
		}

		bankBranch := BankBranch{
			Address:       address,
			BankName:      bankName,
			CountryISO2:   countryISO2,
			IsHeadquarter: isHeadquarter,
			SwiftCode:     swiftCode,
		}

		countryBanks.SwiftCodes = append(countryBanks.SwiftCodes, bankBranch)
		countryBanks.CountryISO2 = countryISO2
		countryBanks.CountryName = countryName

		if !rows.Next() {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &countryBanks, nil
}

func (r *RelationalDB) AddSwiftCodeEntry(b Bank) error {
	query := `INSERT INTO BanksData (address, bankName, countryISO2, countryName, isHeadquarter, swiftCode)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, b.Address, b.BankName, b.CountryISO2, b.CountryName, b.IsHeadquarter, b.SwiftCode)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrSwiftCodeExists
		}
		return err
	}

	return nil
}

func (r *RelationalDB) DeleteSwiftCodeEntry(swiftCode string) error {
	query := `SELECT swiftCode FROM BanksData WHERE swiftCode = $1`
	_, err := r.db.Exec(query, swiftCode)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrSwiftCodeNotFound
	}

	query = `DELETE FROM BanksData WHERE swiftCode = $1`
	_, err = r.db.Exec(query, swiftCode)

	return err
}
