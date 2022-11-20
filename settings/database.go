/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package settings

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

type Database struct {
	SQLite3    string      `json:"sqlite3"`
	PostgreSQL *PostgreSQL `json:"postgresql"`
}

func (d *Database) Build() (*gorm.DB, error) {
	var dialector gorm.Dialector
	if d.SQLite3 != "" {
		dialector = d.buildSQLiteDialector()
	} else if d.PostgreSQL != nil {
		dialector = d.buildPostgreSQL()
	}

	if dialector == nil {
		return nil, errors.New("database is not set")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}

	return db, nil
}

func (d *Database) buildSQLiteDialector() gorm.Dialector {
	return sqlite.Open(d.SQLite3)
}

func (d *Database) buildPostgreSQL() gorm.Dialector {
	if d.PostgreSQL.Port == 0 {
		d.PostgreSQL.Port = 5432
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		getOrDefault(d.PostgreSQL.Host, "localhost"),
		getOrDefault(d.PostgreSQL.User, "postgres"),
		d.PostgreSQL.Password,
		d.PostgreSQL.DbName,
		d.PostgreSQL.Port,
	)
	return postgres.Open(dsn)
}

func getOrDefault(val, default_ string) string {
	if val == "" {
		return default_
	}

	return val
}
