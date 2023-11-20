package usersRepository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		gormModels.User{},
		gormModels.Cookie{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) Create(user httpModels.User) (uint64, error) {
	var usr httpModels.Username

	if err := db.DB.Model(&gormModels.User{}).Create(&gormModels.User{
		Username: user.Username,
		Password: user.Password,
	}).Scan(&usr).Error; err != nil {
		return 0, err
	}

	return usr.ID, nil
}

func (db *Postgres) GetByUsername(username string) (gormModels.User, error) {
	var usr gormModels.User
	if err := db.DB.Model(&gormModels.User{Username: username}).
		Where("username = ?", username).First(&usr).Error; err != nil {
		return gormModels.User{}, err
	}
	return usr, nil
}

func (db *Postgres) GetByID(id uint64) (gormModels.User, error) {
	var usr gormModels.User
	if err := db.DB.Preload("Users").Where("id = ?", id).First(&usr).Error; err != nil {
		return gormModels.User{}, err
	}
	return usr, nil
}

func (db *Postgres) GetUserIDBySessionID(sessionID string) (uint64, error) {
	var usrID gormModels.Cookie
	if err := db.DB.Table(gormModels.SessionTable).Model(&gormModels.User{}).
		Where("value = ?", sessionID).
		Select("user_id").
		Last(&usrID).Error; err != nil {
		return 0, err
	}
	return usrID.UserID, nil
}

func (db *Postgres) GetBySessionID(sessionID string) (gormModels.User, error) {
	var usr gormModels.User
	if err := db.DB.Table("users").Model(&gormModels.User{}).
		Joins("JOIN cookies ON cookies.user_id=users.id").
		Where("cookies.value = ? AND cookies.deleted_at IS NULL", sessionID).
		First(&usr).Error; err != nil {
		return gormModels.User{}, err
	}
	return usr, nil
}
