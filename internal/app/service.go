package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mikekonan/go-countries"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"github.com/pkacprzak5/bic-data-service/pkg/utils"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type BankService struct {
	storage storage.Storage
}

func NewBankService(s storage.Storage) *BankService {
	return &BankService{storage: s}
}

func (s *BankService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /swift-codes/{swift-code}", s.handleGetSwiftCodeDetails)
	router.HandleFunc("GET /swift-codes/country/{countryISO2code}", s.handleGetCountrySwiftCodes)
	router.HandleFunc("POST /swift-codes", s.handleAddSwiftCodeDetails)
	router.HandleFunc("DELETE /swift-codes/{swift-code}", s.handleDeleteSwiftCode)
}

func (s *BankService) handleGetSwiftCodeDetails(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.PathValue("swift-code")
	if swiftCode == "" {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "swift-code is empty"})
		return
	}

	bank, err := s.storage.GetSwiftCodeDetails(swiftCode)
	if err != nil && errors.Is(err, storage.ErrSwiftCodeNotFound) {
		utils.WriteJSON(w, http.StatusNotFound, storage.Response{Message: err.Error()})
		return
	} else if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, storage.Response{Message: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, bank)
}

func (s *BankService) handleGetCountrySwiftCodes(w http.ResponseWriter, r *http.Request) {
	countryISO2code := r.PathValue("country-iso-2code")
	if countryISO2code == "" {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "country-iso-2code is empty"})
		return
	}

	if countryISO2code != strings.ToUpper(countryISO2code) || !isValidISO2(countryISO2code) {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "country-iso-2code is invalid"})
		return
	}

	swiftCodes, err := s.storage.GetSwiftCodesForCountry(countryISO2code)
	if err != nil && errors.Is(err, storage.ErrISO2CodeNotFound) {
		utils.WriteJSON(w, http.StatusNotFound, storage.Response{Message: err.Error()})
		return
	} else if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, storage.Response{Message: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, swiftCodes)
}

func (s *BankService) handleAddSwiftCodeDetails(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "Error reading request body"})
		return
	}

	defer r.Body.Close()

	var bank storage.Bank
	err = json.Unmarshal(body, &bank)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "Error parsing request body"})
		return
	}

	if err := validateBankData(bank); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: err.Error()})
		return
	}

	err = s.storage.AddSwiftCodeEntry(bank)
	if err != nil && errors.Is(err, storage.ErrSwiftCodeExists) {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: err.Error()})
		return
	} else if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, storage.Response{Message: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK,
		storage.Response{Message: fmt.Sprintf("Successfully added bank with swift code %s", bank.SwiftCode)})
}

func (s *BankService) handleDeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.PathValue("swift-code")
	if swiftCode == "" {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "swift-code is empty"})
		return
	}

	err := s.storage.DeleteSwiftCodeEntry(swiftCode)
	if err != nil && errors.Is(err, storage.ErrSwiftCodeNotFound) {
		utils.WriteJSON(w, http.StatusNotFound, storage.Response{Message: err.Error()})
		return
	} else if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, storage.Response{Message: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK,
		storage.Response{Message: fmt.Sprintf("Bank with swift code: %s has been deleted", swiftCode)})
}

func isValidISO2(code string) bool {
	_, flag := country.ByAlpha2CodeStr(code)
	return flag
}

func validateBankData(b storage.Bank) error {
	if strings.TrimSpace(b.Address) == "" {
		return errors.New("address is required")
	}

	if strings.TrimSpace(b.BankName) == "" {
		return errors.New("bankName is required")
	}

	if b.CountryISO2 == "" {
		return errors.New("countryISO2 is required")
	}

	if !isValidISO2(b.CountryISO2) {
		return errors.New("countryISO2 is invalid")
	}

	if b.CountryName == "" {
		return errors.New("countryName is required")
	}

	b.CountryName = strings.ToUpper(b.CountryName)
	if b.CountryName != iso2CodeToCountry(b.CountryISO2) {
		return errors.New("countryName does not match ISO2 code")
	}

	if b.IsHeadquarter == nil {
		return errors.New("isHeadquarter is required")
	}

	if b.SwiftCode == "" {
		return errors.New("swiftCode is required")
	}

	if !isValidSWIFT(b.SwiftCode) {
		return errors.New("swiftCode is invalid")
	}

	if strings.HasSuffix(b.SwiftCode, "XXX") && !*b.IsHeadquarter {
		return errors.New("swiftCode indicates bank's headquarter")
	}

	if *b.IsHeadquarter && !strings.HasSuffix(b.SwiftCode, "XXX") {
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

func isValidSWIFT(s string) bool {
	regex := `^[A-Za-z]{4}[A-Za-z]{2}[A-Za-z0-9]{2}([A-Za-z0-9]{3})?$`
	matched, _ := regexp.MatchString(regex, s)
	return matched
}
