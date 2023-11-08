package models

import "gorm.io/gorm"

type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	gorm.Model
	ID       uint64 `json:"userID"   gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Username struct {
	ID       uint64
	Username string
}

type UsersInfo struct {
	UserID   uint64
	Password string
}

type UserID struct {
	ID uint64 `json:"userID"`
}
