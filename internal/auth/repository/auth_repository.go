package authRepository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/models"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(models.Cookie{})

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateSession(cookie models.Cookie) (string, error) {
	res := db.DB.Table(models.SessionTable).Create(&cookie)
	if res.Error != nil {
		return "", res.Error
	}
	return cookie.Value, nil
}

func (db *Postgres) DeleteBySessionID(sessionID string) error {
	if err := db.DB.Delete(&models.Cookie{}, "value = ?", sessionID).
		Error; err != nil {
		return err
	}
	return nil
}
