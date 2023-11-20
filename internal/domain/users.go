package domain

import httpModels "github.com/cyber_bed/internal/models/http"

type UsersUsecase interface {
	CreateUser(user httpModels.User) (uint64, error)

	GetUserIDBySessionID(sessionID string) (uint64, error)
	GetBySessionID(sessionID string) (httpModels.User, error)
	GetByUsername(username string) (httpModels.User, error)
	GetByID(userID uint64) (httpModels.User, error)
}
