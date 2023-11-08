package models

import (
	"time"

	"gorm.io/gorm"
)

type Cookie struct {
	gorm.Model
	UserID     uint64    `json:"value"       gorm:"foreignKey"`
	Value      string    `json:"userID"`
	ExpireDate time.Time `json:"expire_date"`
}
