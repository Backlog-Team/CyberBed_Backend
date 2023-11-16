package usersRepository

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

	db.AutoMigrate(
		models.User{},
		models.Cookie{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) Create(user models.User) (uint64, error) {
	var usr models.Username

	if err := db.DB.Model(&models.User{}).Create(&user).Scan(&usr).Error; err != nil {
		return 0, err
	}

	return usr.ID, nil
}

func (db *Postgres) GetByUsername(username string) (models.User, error) {
	var usr models.User
	if err := db.DB.Model(&models.User{Username: username}).
		Where("username = ?", username).First(&usr).Error; err != nil {
		return models.User{}, err
	}
	return usr, nil
}

func (db *Postgres) GetByID(id uint64) (models.User, error) {
	var usr models.User
	if err := db.DB.Preload("Users").Where("id = ?", id).First(&usr).Error; err != nil {
		return models.User{}, err
	}
	return usr, nil
}

func (db *Postgres) GetUserIDBySessionID(sessionID string) (uint64, error) {
	var usrID models.Cookie
	if err := db.DB.Table(models.SessionTable).Model(&models.User{}).
		Where("value = ?", sessionID).
		Select("user_id").
		Last(&usrID).Error; err != nil {
		return 0, err
	}
	return usrID.UserID, nil
}

func (db *Postgres) GetBySessionID(sessionID string) (models.User, error) {
	var usr models.User
	if err := db.DB.Table("users").Model(&models.User{}).
		Joins("JOIN cookies ON cookies.user_id=users.id").
		Where("cookies.value = ? AND cookies.deleted_at IS NULL", sessionID).
		First(&usr).Error; err != nil {
		return models.User{}, err
	}
	return usr, nil
}
