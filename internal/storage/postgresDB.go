package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgreSQLStorage struct {
	Db *sql.DB
}

// Helper function that returns errors (for testing)
func newPostgreSQLStorage(config PostgresConfig) (*PostgreSQLStorage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		config.Host, config.DB_Port, config.User, config.Password, config.Database,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return &PostgreSQLStorage{Db: db}, nil
}

// Public function that exits on error (for production)
func NewPostgreSQLStorage(config PostgresConfig) (*PostgreSQLStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config.Host, config.DB_Port, config.User, config.Password, config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to connect to database: %v", err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to ping database: %v", err))
		return nil, err
	}

	log.Println("Successfully connected to database")

	return &PostgreSQLStorage{Db: db}, nil
}

func (s *PostgreSQLStorage) Init() (*sql.DB, error) {
	_, err := s.Db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;")
	if err != nil {
		return nil, fmt.Errorf("failed to create pg_trgm extension: %v", err)
	}
	if err := s.createBankTable(); err != nil {
		return nil, err
	}
	return s.Db, nil
}

func (s *PostgreSQLStorage) createBankTable() error {
	_, err := s.Db.Exec(`
			CREATE TABLE IF NOT EXISTS BanksData (
			    address TEXT NOT NULL,
			    bankName TEXT NOT NULL,
			    isHeadquarter BOOLEAN NOT NULL,
			    countryName TEXT NOT NULL,
			    countryISO2 CHAR(2) NOT NULL,
			    swiftCode TEXT NOT NULL UNIQUE,
			    PRIMARY KEY (swiftCode));

			CREATE INDEX IF NOT EXISTS idx_countryISO2 ON BanksData (countryISO2);

			CREATE INDEX IF NOT EXISTS idx_swiftCode_pattern ON BanksData USING gin (swiftCode gin_trgm_ops);
`)
	return err
}
