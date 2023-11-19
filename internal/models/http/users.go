package httpModels

import gormModels "github.com/cyber_bed/internal/models/gorm"

type AuthUser struct {
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

type User struct {
	ID       uint64 `json:"userID"   gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserGormToHttp(user gormModels.User) User {
	return User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
}
