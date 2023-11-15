package domain

import (
	"github.com/cyber_bed/internal/models"
)

type AuthRepository interface {
	CreateSession(cookie models.Cookie) (string, error)
	DeleteBySessionID(sessionID string) error
}

type UsersRepository interface {
	Create(user models.User) (uint64, error)

	GetUserIDBySessionID(sessionID string) (uint64, error)
	GetByUsername(username string) (models.User, error)
	GetByID(id uint64) (models.User, error)
	GetBySessionID(sessionID string) (models.User, error)
}

type PlantsRepository interface {
	CreateUserPlantsRelations(userID uint64, plantID []int64) error
	AddUserPlantsRelations(userID uint64, plantsID []int64) error
	GetPlantsByUserID(userID uint64) (models.UserPlants, error)
	UpdateUserPlantsRelation(relation models.UserPlants) error
	GetPlantByID(plantID uint64) (models.XiaomiPlant, error)
	GetByPlantName(plantName string) ([]models.XiaomiPlant, error)
	GetPlantsPage(pageNum uint64) ([]models.XiaomiPlant, error)
}

type FoldersRepository interface {
	CreateFolder(folder models.Folder) (uint64, error)
	GetFolders(userID uint64) ([]models.Folder, error)
	GetFolder(id uint64) (models.Folder, error)
	DeleteFolder(id uint64) error
	GetPlantsID(folderID uint64) ([]uint64, error)
}
