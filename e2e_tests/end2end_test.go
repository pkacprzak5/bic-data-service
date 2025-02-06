//go:build end2end

package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

type IntegrationTestSuite struct {
	suite.Suite
	dbConfig  storage.PostgresConfig
	serverCmd *exec.Cmd
	serverURL string
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.serverURL = fmt.Sprintf("http://localhost:%s", storage.GetEnv("PORT", "8081"))
	initTestEnvs()
	s.dbConfig = storage.PostgresConfig{
		DB_Port:  storage.GetEnv("DB_PORT", "5432"),
		User:     storage.GetEnv("DB_USER", "test_user"),
		Password: storage.GetEnv("DB_PASSWORD", "Test@1234"),
		Host:     storage.GetEnv("DB_HOST", "localhost"),
		Database: storage.GetEnv("DB_NAME", "testdatabase"),
	}
}

func initTestEnvs() {
	if _, exists := os.LookupEnv("PORT"); !exists {
		os.Setenv("PORT", "8081")
	}
	if _, exists := os.LookupEnv("DB_PORT"); !exists {
		os.Setenv("DB_PORT", "5432")
	}
	if _, exists := os.LookupEnv("DB_USER"); !exists {
		os.Setenv("DB_USER", "test_user")
	}
	if _, exists := os.LookupEnv("DB_PASSWORD"); !exists {
		os.Setenv("DB_PASSWORD", "Test@1234")
	}
	if _, exists := os.LookupEnv("DB_HOST"); !exists {
		os.Setenv("DB_HOST", "localhost")
	}
	if _, exists := os.LookupEnv("DB_NAME"); !exists {
		os.Setenv("DB_NAME", "testdatabase")
	}
}

func (s *IntegrationTestSuite) BeforeTest(_, _ string) {
	err := buildApp()
	if err != nil {
		s.T().Fatal(err)
		return
	}
	cmd := exec.Command("../bin/api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		s.T().Fatalf("Failed to start server: %v", err)
	}

	s.serverCmd = cmd

	time.Sleep(1 * time.Second)
}

const (
	binaryName = "../bin/api"
	sourceFile = "../cmd/main.go"
)

func buildApp() error {
	cmd := exec.Command("go", "build", "-o", binaryName, sourceFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *IntegrationTestSuite) AfterTest(_, _ string) {
	if s.serverCmd != nil && s.serverCmd.Process != nil {
		if err := s.serverCmd.Process.Signal(os.Interrupt); err != nil {
			fmt.Printf("Failed to kill server: %v", err)
		}
	}

	postgresDB, err := storage.NewPostgreSQLStorage(s.dbConfig)
	if err != nil {
		s.T().Fatal(err)
	}
	defer postgresDB.Db.Close()

	_, err = postgresDB.Db.Exec(`
        SELECT pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = $1
          AND pid <> pg_backend_pid()`, s.dbConfig.Database)
	if err != nil {
		s.T().Errorf("Failed to terminate database connections: %v", err)
	}

	_, err = postgresDB.Db.Exec("DROP TABLE IF EXISTS BanksData CASCADE")
	if err != nil {
		s.T().Errorf("Failed to drop table BanksData: %v", err)
	}
	s.serverCmd = nil
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// test adding, deleting, querying
func (s *IntegrationTestSuite) TestFullAPIFlow() {
	client := &http.Client{
		Timeout: 1 * time.Second, // Adjust timeout as needed
	}

	// Test data
	swiftCode := "TESTPL33XXX"
	data := storage.Bank{
		SwiftCode:     strPtr(swiftCode),
		BankName:      strPtr("Test Bank"),
		CountryISO2:   strPtr("PL"),
		CountryName:   strPtr("Poland"),
		IsHeadquarter: boolPtr(true),
		Address:       strPtr("123 Bank St"),
	}

	body, _ := json.Marshal(data)

	// Test POST (Create)
	postResp, err := client.Post(s.serverURL+"/v1/swift-codes", "application/json", bytes.NewBuffer(body))
	if err != nil {
		s.T().Fatalf("POST request failed: %v", err)
	}
	defer postResp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, postResp.StatusCode)

	// Test GET (Read)
	getURL := fmt.Sprintf("%s/v1/swift-codes/%s", s.serverURL, swiftCode)
	getResp, err := client.Get(getURL)
	if err != nil {
		s.T().Fatalf("GET request failed: %v", err)
	}
	defer getResp.Body.Close()
	var receivedData storage.Bank
	err = json.NewDecoder(getResp.Body).Decode(&receivedData)
	if err != nil {
		s.T().Fatalf("Failed to decode response body: %v", err)
	}
	assert.Equal(s.T(), http.StatusOK, getResp.StatusCode)

	// Test DELETE
	deleteURL := fmt.Sprintf("%s/v1/swift-codes/%s", s.serverURL, swiftCode)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		s.T().Fatalf("Failed to create DELETE request: %v", err)
	}

	deleteResp, err := client.Do(req)
	if err != nil {
		s.T().Fatalf("DELETE request failed: %v", err)
	}
	defer deleteResp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, deleteResp.StatusCode)

	// Test GET (Read) after deleting
	getURL = fmt.Sprintf("%s/v1/swift-codes/%s", s.serverURL, swiftCode)
	getResp, err = client.Get(getURL)
	if err != nil {
		s.T().Fatalf("GET request failed: %v", err)
	}
	defer getResp.Body.Close()
	assert.Equal(s.T(), http.StatusNotFound, getResp.StatusCode)
}

func (s *IntegrationTestSuite) TestQueryByCountryISO2() {
	client := &http.Client{Timeout: 1 * time.Second}

	swiftCode := "TESTPL33XXX"
	data := storage.Bank{
		SwiftCode:     strPtr(swiftCode),
		BankName:      strPtr("Test Bank"),
		CountryISO2:   strPtr("PL"),
		CountryName:   strPtr("Poland"),
		IsHeadquarter: boolPtr(true),
		Address:       strPtr("123 Bank St"),
	}

	body, _ := json.Marshal(data)
	postResp, err := client.Post(s.serverURL+"/v1/swift-codes", "application/json", bytes.NewBuffer(body))
	assert.NoError(s.T(), err)
	postResp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, postResp.StatusCode)

	swiftCode = "TESTPL44XXX"
	data = storage.Bank{
		SwiftCode:     strPtr(swiftCode),
		BankName:      strPtr("Test Bank 2"),
		CountryISO2:   strPtr("PL"),
		CountryName:   strPtr("Poland"),
		IsHeadquarter: boolPtr(true),
		Address:       strPtr("123 Test St"),
	}

	body, _ = json.Marshal(data)
	postResp, err = client.Post(s.serverURL+"/v1/swift-codes", "application/json", bytes.NewBuffer(body))
	assert.NoError(s.T(), err)
	postResp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, postResp.StatusCode)

	getResp, err := client.Get(s.serverURL + "/v1/swift-codes/country/PL")
	assert.NoError(s.T(), err)
	defer getResp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, getResp.StatusCode)

	var country storage.CountryBanks
	err = json.NewDecoder(getResp.Body).Decode(&country)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), country.SwiftCodes, 2)
}

func (s *IntegrationTestSuite) TestQueryByCountryWithNoData() {
	client := &http.Client{Timeout: 1 * time.Second}

	getResp, err := client.Get(s.serverURL + "/v1/swift-codes/country/US")
	assert.NoError(s.T(), err)
	defer getResp.Body.Close()
	assert.Equal(s.T(), http.StatusNotFound, getResp.StatusCode)
}

func (s *IntegrationTestSuite) TestInvalidJSONBody() {
	client := &http.Client{Timeout: 1 * time.Second}

	invalidJSON := `{"swift_code": "TEST", "bank_name": "Test Bank" "country_iso2": "PL"}`
	postResp, err := client.Post(s.serverURL+"/v1/swift-codes", "application/json", bytes.NewBufferString(invalidJSON))
	assert.NoError(s.T(), err)
	defer postResp.Body.Close()
	assert.Equal(s.T(), http.StatusBadRequest, postResp.StatusCode)
}

func (s *IntegrationTestSuite) TestInvalidCountryISO2() {
	client := &http.Client{Timeout: 1 * time.Second}

	getResp, err := client.Get(s.serverURL + "/v1/swift-codes/country/ZZ")
	assert.NoError(s.T(), err)
	defer getResp.Body.Close()
	assert.Equal(s.T(), http.StatusBadRequest, getResp.StatusCode)
}

func (s *IntegrationTestSuite) TestDeleteNonExistentSwiftCode() {
	client := &http.Client{Timeout: 1 * time.Second}

	req, err := http.NewRequest(http.MethodDelete, s.serverURL+"/v1/swift-codes/NON_EXISTENT_CODE", nil)
	assert.NoError(s.T(), err)

	deleteResp, err := client.Do(req)
	assert.NoError(s.T(), err)
	defer deleteResp.Body.Close()
	assert.Equal(s.T(), http.StatusBadRequest, deleteResp.StatusCode)
}
