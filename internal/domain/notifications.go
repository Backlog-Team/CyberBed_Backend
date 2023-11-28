package domain

import (
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type NotificationsUsecase interface {
	CreateNotification(notification httpModels.Notification) (httpModels.Notification, error)
	GetNotificationsByUserID(userID uint64) ([]httpModels.Notification, error)
	GetNotificationsByUserIDAndStatus(
		userID uint64,
		status gormModels.NotificationStatus,
	) ([]gormModels.Notification, error)
	GetNotificationByID(id uint64) (httpModels.Notification, error)
	DeleteNotification(id uint64) error
	UpdateNotificationStatus(
		id uint64,
		status gormModels.NotificationStatus,
	) error
	DeleteNotificationByIDAndStatus(
		id uint64,
		status gormModels.NotificationStatus,
	) error
	UpdatePeriodNotification(notification httpModels.Notification) error
}
