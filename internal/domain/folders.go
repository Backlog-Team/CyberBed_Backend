package domain

import "github.com/cyber_bed/internal/models"

type FoldersUsecase interface {
	CreateFolder(folder models.Folder) (uint64, error)
	GetFolderByID(id uint64) (models.FolderHttp, error)
	GetFoldersByUserID(userID uint64) ([]models.FolderHttp, error)
	GetPlantsFromFolder(folderID uint64) ([]models.XiaomiPlant, error)
	DeleteFolderByID(id uint64) error
	AddPlantToFolder(folderID, plantID uint64) error
	DeletePlantFromFolder(folderID, plantID uint64) error
}
