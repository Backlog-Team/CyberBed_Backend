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
}

type FoldersRepository interface {
	CreateFolder(folder httpModels.Folder) (uint64, error)
	GetFolders(userID uint64) ([]gormModels.Folder, error)
	GetFolder(id uint64) (gormModels.Folder, error)
	DeleteFolder(id uint64) error
	GetPlantsID(folderID uint64) ([]uint64, error)
	AddPlantToFolder(folderID, plantID uint64) error
	UpdateFolderPlant(folderID, plantID uint64) error
}
