package settings

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	SQLite3File string `json:"sqlite3"`
}

func (d *Database) Create() (*gorm.DB, error) {
	var dialector gorm.Dialector
	if d.SQLite3File != "" {
		dialector = sqlite.Open(d.SQLite3File)
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
