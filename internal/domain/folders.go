package domain

import httpModels "github.com/cyber_bed/internal/models/http"

type FoldersUsecase interface {
	CreateFolder(folder httpModels.Folder) (uint64, error)
	GetFolderByID(id uint64) (httpModels.Folder, error)
	GetFoldersByUserID(userID uint64) ([]httpModels.Folder, error)
	GetPlantsFromFolder(folderID uint64) ([]httpModels.XiaomiPlant, error)
	DeleteFolderByID(id uint64) error
	AddPlantToFolder(folderID, plantID uint64) error
	DeletePlantFromFolder(folderID, plantID uint64) error
	GetFolderByPlantAndUserID(plantID, userID uint64) (map[httpModels.Folder]map[uint64]bool, error)
	CreateChannel(folderID, plantID, channelID uint64) (uint64, error)
	GetChannelByFolderPlantID(folderID, plantID uint64) (uint64, error)
}
