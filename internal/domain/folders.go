package domain

import "github.com/cyber_bed/internal/models"

type FoldersUsecase interface {
	CreateFolder(folder models.Folder) (uint64, error)
	GetFolderByID(id uint64) (models.Folder, error)
	GetFoldersByUserID(userID uint64) ([]models.Folder, error)
	GetPlantsFromFolder(folderID uint64) ([]models.XiaomiPlant, error)
	DeleteFolderByID(id uint64) error
}
