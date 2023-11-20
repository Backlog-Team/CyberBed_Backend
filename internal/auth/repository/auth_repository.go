package authRepository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormModels "github.com/cyber_bed/internal/models/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(gormModels.Cookie{})

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateSession(cookie gormModels.Cookie) (string, error) {
	res := db.DB.Table(gormModels.SessionTable).Create(&cookie)
	if res.Error != nil {
		return "", res.Error
	}
	return cookie.Value, nil
}

func (db *Postgres) DeleteBySessionID(sessionID string) error {
	if err := db.DB.Delete(&gormModels.Cookie{}, "value = ?", sessionID).
		Error; err != nil {
		return err
	}
	return nil
}
