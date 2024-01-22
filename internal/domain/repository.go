package domain

import (
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type AuthRepository interface {
	CreateSession(cookie gormModels.Cookie) (string, error)
	DeleteBySessionID(sessionID string) error
}

type UsersRepository interface {
	Create(user httpModels.User) (uint64, error)

	GetUserIDBySessionID(sessionID string) (uint64, error)
	GetByUsername(username string) (gormModels.User, error)
	GetByID(id uint64) (gormModels.User, error)
	GetBySessionID(sessionID string) (gormModels.User, error)
}

type PlantsRepository interface {
	CreateUserPlantsRelations(userID uint64, plantID []int64) error
	AddUserPlantsRelations(userID uint64, plantsID []int64) error
	GetPlantsByUserID(userID uint64) (gormModels.UserPlants, error)
	UpdateUserPlantsRelation(relation gormModels.UserPlants) error
	GetPlantByID(plantID uint64) (gormModels.XiaomiPlant, error)
	GetByPlantName(plantName string) ([]gormModels.XiaomiPlant, error)
	GetPlantsPage(pageNum uint64) ([]gormModels.XiaomiPlant, error)

	CreateCustomPlant(plant httpModels.CustomPlant) (uint64, error)
	UpdateCustomPlant(plant httpModels.CustomPlant) error
	GetCustomPlants(userID uint64) ([]gormModels.CustomPlant, error)
	GetCustomPlant(userID, plantID uint64) (gormModels.CustomPlant, error)
	DeleteCustomPlant(userID, plantID uint64) error

	CreateSavedPlant(userID, plantID uint64) error
	GetSavedPlants(userID uint64) ([]gormModels.SavedPlant, error)
	DeleteSavedPlant(userID, plantID uint64) error
	GetSavedPlantByIDs(userID, plantID uint64) (gormModels.SavedPlant, error)

	CreateChannel(plantID, channelID, userID uint64) (uint64, error)
	GetChannelByUserAndPlantID(userID, plantID uint64) (uint64, error)
	UpdateChannelByUserAndPlantID(userID, plantID, channelID uint64) error
}

type FoldersRepository interface {
	CreateFolder(folder httpModels.Folder) (uint64, error)
	GetFolders(userID uint64) ([]gormModels.Folder, error)
	GetFolder(id uint64) (gormModels.Folder, error)
	GetFolderByNameAndUserID(folderName string, userID uint64) (gormModels.Folder, error)
	DeleteFolder(id uint64) error
	GetPlantsID(folderID uint64) ([]uint64, error)
	AddPlantToFolder(folderID, plantID uint64) error
	UpdateFolderPlant(folderID, plantID uint64) error
	GetFolderByPlantAndUserID(userID, plantID uint64) ([]gormModels.Folder, error)
}

type NotificationsRepository interface {
	CreateNotification(notification gormModels.Notification) (gormModels.Notification, error)
	GetNotificationsByUserID(userID uint64) ([]gormModels.Notification, error)
	GetNotificationsByUserIDAndStatus(
		userID uint64,
		status gormModels.NotificationStatus,
	) ([]gormModels.Notification, error)
	GetNotificationByID(id uint64) (gormModels.Notification, error)
	GetNotificationsByUserPlantID(userID uint64, plantID uint64) (gormModels.Notification, error)
	GetWaitingNotification(userID, plantID uint64) (gormModels.Notification, error)
	UpdateNotificationStatus(id uint64, status gormModels.NotificationStatus) error
	DeleteNotification(id uint64) error
	DeleteNotificationByIDAndStatus(id uint64, status gormModels.NotificationStatus) error
	UpdatePeriodNotification(notification gormModels.Notification) (gormModels.Notification, error)
}
