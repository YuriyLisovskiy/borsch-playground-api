package db

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgreSQLFromEnv() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	portString := os.Getenv("DB_PORT")
	dbPort := 5432
	if portString != "" {
		port, err := strconv.Atoi(portString)
		if err != nil {
			return nil, err
		}

		dbPort = port
	}

	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")
	options := os.Getenv("DB_OPTIONS")
	dialector := buildPostgreSQL(dbPort, host, user, password, dbName, options)
	if dialector == nil {
		return nil, errors.New("database is not set")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}

	return db, nil
}

func buildPostgreSQL(port int, host, user, password, dbName, options string) gorm.Dialector {
	if port == 0 {
		port = 5432
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d %s",
		getOrDefault(host, "localhost"),
		getOrDefault(user, "postgres"),
		password,
		dbName,
		port,
		options,
	)
	return postgres.Open(dsn)
}

func getOrDefault(val, default_ string) string {
	if val == "" {
		return default_
	}

	return val
}
