package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Config holds application configuration
type Config struct {
	Port     string
	DBConfig *DBConfig
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var (
	// DB is the global database connection pool
	DB *sql.DB
)

// Load returns application configuration
func Load() (*Config, error) {
	dbConfig := &DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "gratitude_user"),
		Password: getEnvOrDefault("DB_PASSWORD", "gratitude123"),
		DBName:   getEnvOrDefault("DB_NAME", "gratitude_jar"),
	}

	return &Config{
		Port:     getEnvOrDefault("PORT", ":4000"),
		DBConfig: dbConfig,
	}, nil
}

// InitDB initializes the database connection
func InitDB() error {
	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	dbCfg := cfg.DBConfig
	log.Printf("Connecting to database with host=%s port=%s user=%s dbname=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.DBName)

	// Construct the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName,
	)

	// Open the database connection
	log.Printf("Opening database connection...")
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	log.Printf("Testing database connection...")
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("Successfully connected to database")
	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
