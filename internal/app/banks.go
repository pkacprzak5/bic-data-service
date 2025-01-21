package app

import (
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"net/http"
)

type BankService struct {
	storage storage.Storage
}

func NewBankService(s storage.Storage) *BankService {
	return &BankService{storage: s}
}

func (s *BankService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /swift-codes/{swift-code}", s.GetSwiftCodeDetails)
	router.HandleFunc("GET /swift-codes/country/{countryISO2code}", s.GetCountrySwiftCodes)
	router.HandleFunc("POST /swift-codes", s.AddSwiftCodeDetails)
	router.HandleFunc("DELETE /swift-codes/{swift-code}", s.DeleteSwiftCode)
}

func (s *BankService) GetSwiftCodeDetails(writer http.ResponseWriter, request *http.Request) {
	
}

func (s *BankService) GetCountrySwiftCodes(writer http.ResponseWriter, request *http.Request) {

}

func (s *BankService) AddSwiftCodeDetails(writer http.ResponseWriter, request *http.Request) {

}

func (s *BankService) DeleteSwiftCode(writer http.ResponseWriter, request *http.Request) {

}
