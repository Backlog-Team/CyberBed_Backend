package httpModels

import (
	"time"

	gormModels "github.com/cyber_bed/internal/models/gorm"
)

type Notification struct {
	ID             uint64                        `json:"id"`
	UserID         uint64                        `json:"userID"`
	FolderID       uint64                        `json:"folderID"`
	PlantID        uint64                        `json:"plantID"`
	ExpirationTime string                        `json:"period"`
	Status         gormModels.NotificationStatus `json:"status"`
	TimeLeft       string                        `json:"timeLeft"`
}

func NotificationGormToHttp(notification gormModels.Notification) Notification {
	return Notification{
		ID:             notification.ID,
		UserID:         notification.UserID,
		PlantID:        notification.PlantID,
		FolderID:       notification.FolderID,
		Status:         notification.Status,
		ExpirationTime: notification.Period,
		TimeLeft:       notification.ExpirationTime.Sub(time.Now()).String(),
	}
}
