package migrations

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormModels "github.com/cyber_bed/internal/models/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func newPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB: db,
	}, nil
}

func StartMigration(url, pathToDir string) error {
	entries, err := os.ReadDir(pathToDir)
	filesNum := len(entries)
	if err != nil {
		return err
	}

	db, err := newPostgres(url)
	if err != nil {
		return err
	}

	for plantIndx, e := range entries {
		content, err := os.ReadFile(pathToDir + "/" + e.Name())
		if err != nil {
			return err
		}

		var plantItem gormModels.XiaomiPlant
		json.Unmarshal(content, &plantItem)
		if err = db.createPlant(plantItem); err != nil {
			return err
		} else {
			log.Printf("Written %d of %d plants", plantIndx, filesNum)
		}
	}

	return nil
}

func (db *Postgres) createPlant(plant gormModels.XiaomiPlant) error {
	err := db.DB.Model(&gormModels.XiaomiPlant{}).Create(&plant).Error
	if err != nil {
		return err
	}
	return nil
}
