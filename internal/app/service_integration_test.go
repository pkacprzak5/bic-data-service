//go:build integration

package app

import (
	"encoding/json"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	db        *storage.PostgreSQLStorage
	service   *BankService
	testSwift string
}

func (s *IntegrationTestSuite) SetupSuite() {
	config := storage.PostgresConfig{
		DB_Port:  storage.GetEnv("DB_PORT", "5432"),
		User:     storage.GetEnv("DB_USER", "test_user"),
		Password: storage.GetEnv("DB_PASSWORD", "Test@1234"),
		Host:     storage.GetEnv("DB_HOST", "localhost"),
		Database: storage.GetEnv("DB_NAME", "testdatabase"),
	}

	var err error
	s.db, err = storage.NewPostgreSQLStorage(config)
	if err != nil {
		s.T().Fatalf("Failed to connect to test DB: %v", err)
	}

	_, err = s.db.Init()
	if err != nil {
		s.T().Fatalf("Failed to init DB: %v", err)
	}

	s.service = NewBankService(storage.NewRelationalDB(s.db.Db))
	s.testSwift = "TESTPL44XXX"
}

func (s *IntegrationTestSuite) TearDownSuite() {
	_ = s.db.Db.Close()
}

func (s *IntegrationTestSuite) BeforeTest(_, _ string) {
	_, _ = s.db.Init()
	err := s.db.Db.QueryRow(
		`INSERT INTO BanksData 
		(address, bankName, countryISO2, countryName, isHeadquarter, swiftCode) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"Integration Test Address",
		"Integration Bank",
		"PL",
		"POLAND",
		true,
		s.testSwift,
	).Err()

	if err != nil {
		s.T().Fatalf("Failed to seed data: %v", err)
	}
}

func (s *IntegrationTestSuite) AfterTest(_, _ string) {
	_, _ = s.db.Db.Exec("DROP TABLE IF EXISTS BanksData CASCADE")
}

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestGetSwiftCodeDetails() {
	req := httptest.NewRequest(http.MethodGet, "/swift-codes/"+s.testSwift, nil)
	req = setPathVars(req, map[string]string{"swiftCode": s.testSwift})
	rec := httptest.NewRecorder()

	s.service.handleGetSwiftCodeDetails(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var bank storage.Bank
	assert.NoError(s.T(), json.NewDecoder(res.Body).Decode(&bank))
	assert.Equal(s.T(), s.testSwift, *bank.SwiftCode)
}

func (s *IntegrationTestSuite) TestFullFlow() {
	addReq := httptest.NewRequest(http.MethodPost, "/swift-codes", strings.NewReader(`{
		"address": "New Address",
		"bankName": "New Bank",
		"countryISO2": "PL",
		"countryName": "POLAND",
		"isHeadquarter": true,
		"swiftCode": "TESTPL34XXX"
	}`))

	// adding test
	addRec := httptest.NewRecorder()
	s.service.handleAddSwiftCodeDetails(addRec, addReq)
	assert.Equal(s.T(), http.StatusOK, addRec.Result().StatusCode)

	// getting test
	getReq := httptest.NewRequest(http.MethodGet, "/swift-codes/TESTPL34XXX", nil)
	getReq = setPathVars(getReq, map[string]string{"swiftCode": "TESTPL34XXX"})
	getRec := httptest.NewRecorder()
	s.service.handleGetSwiftCodeDetails(getRec, getReq)
	assert.Equal(s.T(), http.StatusOK, getRec.Result().StatusCode)

	// deleting test
	deleteReq := httptest.NewRequest(http.MethodDelete, "/swift-codes/TESTPL34XXX", nil)
	deleteReq = setPathVars(deleteReq, map[string]string{"swiftCode": "TESTPL34XXX"})
	deleteRec := httptest.NewRecorder()
	s.service.handleDeleteSwiftCode(deleteRec, deleteReq)
	assert.Equal(s.T(), http.StatusOK, deleteRec.Result().StatusCode)
}

func (s *IntegrationTestSuite) TestCountryCodes() {
	_, err := s.db.Db.Exec(`INSERT INTO BanksData 
		(address, bankName, countryISO2, countryName, isHeadquarter, swiftCode) 
		VALUES 
		('Addr1', 'Bank1', 'IT', 'Italy', true, 'TESTIT33XXX'),
		('Addr2', 'Bank2', 'IT', 'Italy', false, 'TESTIT33AAA')`)
	if err != nil {
		s.T().Fatalf("Failed to seed country data: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/swift-codes/country/IT", nil)
	req = setPathVars(req, map[string]string{"countryISO2code": "IT"})
	rec := httptest.NewRecorder()

	s.service.handleGetCountrySwiftCodes(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var country storage.CountryBanks
	assert.NoError(s.T(), json.NewDecoder(res.Body).Decode(&country))
	assert.Len(s.T(), country.SwiftCodes, 2)
}

func (s *IntegrationTestSuite) TestCountryCodeWithNoData() {
	req := httptest.NewRequest(http.MethodGet, "/swift-codes/country/IT", nil)
	req = setPathVars(req, map[string]string{"countryISO2code": "IT"})
	rec := httptest.NewRecorder()

	s.service.handleGetCountrySwiftCodes(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(s.T(), http.StatusNotFound, res.StatusCode)
}
