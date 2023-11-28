package notificationsRepository

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

	db.AutoMigrate(&gormModels.Notification{})

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateNotification(notification gormModels.Notification) (uint64, error) {
	var resRow gormModels.Notification
	if err := db.DB.Create(&notification).Scan(&resRow).Error; err != nil {
		return 0, err
	}
	return resRow.ID, nil
}

func (db *Postgres) GetNotificationsByUserID(userID uint64) ([]gormModels.Notification, error) {
	var userNotifications []gormModels.Notification
	if err := db.DB.Model(&gormModels.Notification{
		UserID: userID,
	}).Find(&userNotifications).Error; err != nil {
		return []gormModels.Notification{}, err
	}
	return userNotifications, nil
}

func (db *Postgres) GetNotificationsByUserIDAndStatus(
	userID uint64,
	status gormModels.NotificationStatus,
) ([]gormModels.Notification, error) {
	var userNotifications []gormModels.Notification
	if err := db.DB.Model(&gormModels.Notification{}).
		Where("user_id = ? AND status = ?", userID, status).
		Find(&userNotifications).Error; err != nil {
		return []gormModels.Notification{}, err
	}
	return userNotifications, nil
}

func (db *Postgres) UpdateNotificationStatus(
	id uint64,
	status gormModels.NotificationStatus,
) error {
	if err := db.DB.Model(&gormModels.Notification{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetNotificationByID(id uint64) (gormModels.Notification, error) {
	var userNotification gormModels.Notification
	if err := db.DB.Model(&gormModels.Notification{ID: id}).
		First(&userNotification).Error; err != nil {
		return gormModels.Notification{}, err
	}
	return userNotification, nil
}

func (db *Postgres) DeleteNotification(id uint64) error {
	if err := db.DB.Where("id = ?", id).
		Delete(&gormModels.Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) DeleteNotificationByIDAndStatus(
	id uint64,
	status gormModels.NotificationStatus,
) error {
	if err := db.DB.Where("user_id = ? AND status = ?", id, status).
		Delete(&gormModels.Notification{}).Error; err != nil {
		return err
	}
	return nil
}
