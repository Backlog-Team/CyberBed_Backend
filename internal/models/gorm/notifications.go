package gormModels

import (
	"time"

	"gorm.io/gorm"
)

type NotificationStatus string

const (
	NotificationStatusWaiting = "wait"     // waiting for the end for notification
	NotificationStatusDone    = "done"     // done for sending notification
	NotificationStatusFinish  = "archived" // archived for notifications
)

type Notification struct {
	gorm.Model
	ID             uint64
	UserID         uint64
	PlantID        uint64
	ExpirationTime time.Time
	Status         NotificationStatus `gorm:"default:wait"`
}

func (n *Notification) AfterFind(tx *gorm.DB) (err error) {
	if n.Status == NotificationStatusWaiting {
		if time.Now().After(n.ExpirationTime) {
			tx.Model(&Notification{}).Where("id = ?", n.ID).Update("status", NotificationStatusDone)
		}
	}
	return
}
