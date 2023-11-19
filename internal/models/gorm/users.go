package gormModels

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint64 `json:"userID"   gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}
