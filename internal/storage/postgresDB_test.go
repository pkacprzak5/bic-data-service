package storage

import (
	_ "database/sql"
	_ "github.com/lib/pq"
	"os"
	"os/exec"
	"testing"
)

func getValidTestConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "db",
		DB_Port:  "5432",
		User:     "test_user",
		Password: "Test@1234",
		Database: "testdatabase",
	}
}

func getInvalidTestConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "invalid-host",
		DB_Port:  "1234",
		User:     "invalid-user",
		Password: "wrong-password",
		Database: "invalid-db",
	}
}

// TestNewPostgreSQLStorage_Success tests successful connection
func TestNewPostgreSQLStorage_Success(t *testing.T) {
	config := getValidTestConfig()
	storage, err := newPostgreSQLStorage(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer storage.Db.Close()

	if storage.Db == nil {
		t.Error("Database connection is nil")
	}
}

// TestNewPostgreSQLStorage_InvalidConfig tests connection failure
func TestNewPostgreSQLStorage_InvalidConfig(t *testing.T) {
	config := getInvalidTestConfig()
	_, err := newPostgreSQLStorage(config)
	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

// TestNewPostgreSQLStorage_FatalExit tests exit behavior (using a subprocess)
func TestNewPostgreSQLStorage_FatalExit(t *testing.T) {
	if os.Getenv("TEST_FATAL_EXIT") == "1" {
		config := getInvalidTestConfig()
		NewPostgreSQLStorage(config) // Should call log.Fatalf
		return
	}

	// Execute the test in a subprocess to avoid terminating the test runner
	cmd := exec.Command(os.Args[0], "-test.run=TestNewPostgreSQLStorage_FatalExit")
	cmd.Env = append(os.Environ(), "TEST_FATAL_EXIT=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatal("Expected process to exit with error, got success")
	}
}

// TestInit_CreatesSchema tests schema initialization (tables, indexes, extensions).
func TestInit_CreatesSchema(t *testing.T) {
	config := getValidTestConfig()
	storage, err := NewPostgreSQLStorage(config)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer storage.Db.Close()

	defer func() {
		_, err := storage.Db.Exec("DROP TABLE IF EXISTS BanksData CASCADE")
		if err != nil {
			t.Errorf("Failed to clean up tables: %v", err)
		}
	}()

	db, err := storage.Init()
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var extName string
	err = db.QueryRow("SELECT extname FROM pg_extension WHERE extname = 'pg_trgm'").Scan(&extName)
	if err != nil {
		t.Fatalf("pg_trgm extension not found: %v", err)
	}

	_, err = db.Exec("SELECT address, bankname, isheadquarter, countryname, countryiso2, swiftcode FROM banksdata LIMIT 0")
	if err != nil {
		t.Fatalf("BanksData table schema mismatch: %v", err)
	}

	indexes := []string{"idx_countryiso2", "idx_swiftcode_pattern"}
	for _, idx := range indexes {
		var exists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM pg_indexes WHERE indexname = $1)",
			idx,
		).Scan(&exists)
		if err != nil || !exists {
			t.Errorf("Index %s not found or error: %v", idx, err)
		}
	}
}

// TestInit_IsIdempotent tests that multiple Init calls don't cause errors.
func TestInit_IsIdempotent(t *testing.T) {
	config := getValidTestConfig()
	storage, err := NewPostgreSQLStorage(config)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer storage.Db.Close()

	defer func() {
		_, err := storage.Db.Exec("DROP TABLE IF EXISTS BanksData CASCADE")
		if err != nil {
			t.Errorf("Failed to clean up tables: %v", err)
		}
	}()

	if _, err = storage.Init(); err != nil {
		t.Fatalf("First Init failed: %v", err)
	}

	if _, err = storage.Init(); err != nil {
		t.Fatalf("Second Init failed: %v", err)
	}
}

// TestBanksData_UniqueSwiftCode tests the primary key/unique constraint on swiftCode.
func TestBanksData_UniqueSwiftCode(t *testing.T) {
	config := getValidTestConfig()
	storage, err := NewPostgreSQLStorage(config)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer storage.Db.Close()

	defer func() {
		_, err := storage.Db.Exec("DROP TABLE IF EXISTS BanksData CASCADE")
		if err != nil {
			t.Errorf("Failed to clean up tables: %v", err)
		}
	}()

	if _, err = storage.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	_, err = storage.Db.Exec(
		`INSERT INTO BanksData 
			(address, bankname, isheadquarter, countryname, countryiso2, swiftcode)
		VALUES 
			($1, $2, $3, $4, $5, $6)`,
		"Address 1", "Bank 1", true, "Country 1", "C1", "SWIFT1",
	)
	if err != nil {
		t.Fatalf("Failed to insert first record: %v", err)
	}

	_, err = storage.Db.Exec(
		`INSERT INTO BanksData 
			(address, bankname, isheadquarter, countryname, countryiso2, swiftcode)
		VALUES 
			($1, $2, $3, $4, $5, $6)`,
		"Address 2", "Bank 2", false, "Country 2", "C2", "SWIFT1",
	)
	if err == nil {
		t.Error("Expected error for duplicate swiftCode, got nil")
	}
}
