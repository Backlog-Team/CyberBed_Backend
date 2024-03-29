package httpModels

import gormModels "github.com/cyber_bed/internal/models/gorm"

type Folder struct {
	ID         uint64 `json:"ID"`
	UserID     uint64 `json:"userID"`
	FolderName string `json:"folderName"`
	PlantsNum  uint64 `json:"plantsNum"`
	IsDefault  bool   `json:"is_default"`
}

func FolderGormToHttp(f gormModels.Folder) Folder {
	return Folder{
		ID:         f.ID,
		UserID:     f.UserID,
		FolderName: f.FolderName,
		PlantsNum:  f.PlantsNum,
		IsDefault:  f.IsDefalut,
	}
}
