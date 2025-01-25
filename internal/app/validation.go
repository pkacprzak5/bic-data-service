package app

import (
	"errors"
	country "github.com/mikekonan/go-countries"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"regexp"
	"strings"
)

func isValidISO2(code string) bool {
	_, flag := country.ByAlpha2CodeStr(code)
	return flag
}

func validateBankData(b storage.Bank) error {
	if b.Address == nil {
		return errors.New("address is required")
	}

	if b.BankName == nil || strings.TrimSpace(*b.BankName) == "" {
		return errors.New("bankName is required")
	}

	if b.CountryISO2 == nil {
		return errors.New("countryISO2 is required")
	}

	if b.CountryName == nil {
		return errors.New("countryName is required")
	}

	if b.IsHeadquarter == nil {
		return errors.New("isHeadquarter is required")
	}

	if b.SwiftCode == nil {
		return errors.New("swiftCode is required")
	}

	if !isValidISO2(*b.CountryISO2) {
		return errors.New("countryISO2 is invalid")
	}

	*b.CountryName = strings.ToUpper(*b.CountryName)
	if *b.CountryName != iso2CodeToCountry(*b.CountryISO2) {
		return errors.New("countryName does not match ISO2 code")
	}

	if !isValidSWIFT(*b.SwiftCode, *b.CountryISO2) {
		return errors.New("swiftCode is invalid")
	}

	if strings.HasSuffix(*b.SwiftCode, "XXX") && !*b.IsHeadquarter {
		return errors.New("swiftCode indicates bank's headquarter")
	}

	if *b.IsHeadquarter && !strings.HasSuffix(*b.SwiftCode, "XXX") {
		return errors.New("swiftCode indicates bank's branch")
	}

	return nil
}

func iso2CodeToCountry(code string) string {
	countryByISO2Code, ok := country.ByAlpha2CodeStr(code)
	if ok {
		name := countryByISO2Code.NameStr()
		name = strings.ToUpper(name)
		return name
	}
	return ""
}

func isValidSWIFT(s, iso2Code string) bool {
	if s[4:6] != iso2Code { // it should match countryISO2Code of bank location
		return false
	}
	regex := `^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$`
	matched, _ := regexp.MatchString(regex, s)
	return matched
}
