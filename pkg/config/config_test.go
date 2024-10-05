package config

import (
	"testing"
)

func TestDatabaseConfig_ToMigrationUri(t *testing.T) {
	// Test case
	dbConfig := DatabaseConfig{
		Username: "testuser",
		Password: "testpass",
		Host:     "testhost",
		Port:     "5432",
		Database: "testdb",
		SSLMode:  "disable",
	}

	expected := "pgx5://testuser:testpass@testhost:5432/testdb?sslmode=disable"
	result := dbConfig.ToMigrationUri()

	if result != expected {
		t.Errorf("ToMigrationUri() = %v, want %v", result, expected)
	}
}

func TestDatabaseConfig_ToDbConnectionUri(t *testing.T) {
	// Test case
	dbConfig := DatabaseConfig{
		Username: "testuser",
		Password: "testpass",
		Host:     "testhost",
		Port:     "5432",
		Database: "testdb",
		SSLMode:  "require",
	}

	expected := "postgres://testuser:testpass@testhost:5432/testdb?sslmode=require"
	result := dbConfig.ToDbConnectionUri()

	if result != expected {
		t.Errorf("ToDbConnectionUri() = %v, want %v", result, expected)
	}
}
