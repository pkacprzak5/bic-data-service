package storage

type Response struct {
	Message string `json:"message"`
}

type Bank struct {
	Address       string       `json:"address"`
	BankName      string       `json:"bankName"`
	CountryISO2   string       `json:"countryISO2"`
	CountryName   string       `json:"countryName"`
	IsHeadquarter *bool        `json:"isHeadquarter"`
	SwiftCode     string       `json:"swiftCode"`
	Branches      []BankBranch `json:"branches"`
}

type BankBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type CountryBanks struct {
	CountryISO2 string       `json:"countryISO2"`
	CountryName string       `json:"countryName"`
	SwiftCodes  []BankBranch `json:"swiftCode"`
}
