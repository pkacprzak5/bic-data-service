package storage

type Storage interface {
	GetSwiftCodeDetails(swiftCode string) (Bank, error)

	GetSwiftCodesForCountry(iso2Code string) (CountryBanks, error)

	AddSwiftCodeEntry(b Bank) error

	DeleteSwiftCodeEntry(swiftCode string) error
}
