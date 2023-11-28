package httpModels

import (
	gormModels "github.com/cyber_bed/internal/models/gorm"
)

type Notification struct {
	ID             uint64 `json:"id"`
	UserID         uint64 `json:"userID"`
	PlantID        uint64 `json:"plantID"`
	ExpirationTime string `json:"expTime"`
}

type NotificationCancel struct {
	CancelAll bool                          `json:"cancelStatus"`
	CancelID  bool                          `json:"cancelID"`
	Status    gormModels.NotificationStatus `json:"status"`
}

func NotificationGormToHttp(notification gormModels.Notification) Notification {
	return Notification{
		ID:      notification.ID,
		UserID:  notification.UserID,
		PlantID: notification.PlantID,
	}
}
