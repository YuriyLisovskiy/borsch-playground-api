package migrations

import (
	"github.com/YuriyLisovskiy/borsch-playground-service/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.JobOutputRowDbModel{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.JobDbModel{}); err != nil {
		return err
	}

	return nil
}
