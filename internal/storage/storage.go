package storage

import "errors"

type Storage interface {
	GetSwiftCodeDetails(swiftCode string) (Bank, error)

	GetSwiftCodesForCountry(iso2Code string) (CountryBanks, error)

	AddSwiftCodeEntry(b Bank) error

	DeleteSwiftCodeEntry(swiftCode string) error
}

var ErrSwiftCodeNotFound = errors.New("Given Swift Code not found")
var ErrISO2CodeNotFound = errors.New("Country with given ISO2 Code does not have any swift codes")
var ErrSwiftCodeExists = errors.New("Given Swift Code already exists")
