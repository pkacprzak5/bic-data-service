package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"github.com/pkacprzak5/bic-data-service/pkg/utils"
	"io"
	"net/http"
	"strings"
)

type BankService struct {
	storage storage.Storage
}

func NewBankService(s storage.Storage) *BankService {
	return &BankService{storage: s}
}

func (s *BankService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /swift-codes/{swiftCode}", s.handleGetSwiftCodeDetails)
	router.HandleFunc("GET /swift-codes/country/{countryISO2code}", s.handleGetCountrySwiftCodes)
	router.HandleFunc("POST /swift-codes", s.handleAddSwiftCodeDetails)
	router.HandleFunc("DELETE /swift-codes/{swiftCode}", s.handleDeleteSwiftCode)
}

func (s *BankService) handleGetSwiftCodeDetails(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.PathValue("swiftCode")
	if swiftCode == "" {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "swiftCode is empty"})
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
	countryISO2code := r.PathValue("countryISO2code")
	if countryISO2code == "" {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "countryISO2code not found in path"})
		return
	}

	if countryISO2code != strings.ToUpper(countryISO2code) || !isValidISO2(countryISO2code) {
		utils.WriteJSON(w, http.StatusBadRequest, storage.Response{Message: "countryISO2code is invalid"})
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

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError,
				storage.Response{Message: "Error closing body: " + err.Error()})
			return
		}
	}(r.Body)

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
		storage.Response{Message: fmt.Sprintf("Successfully added bank with swift code %s", *bank.SwiftCode)})
}

func (s *BankService) handleDeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.PathValue("swiftCode")
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
