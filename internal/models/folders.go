package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	ID             uint64
	UserID         uint64
	FolderName     string
	PlantsNum      uint64
	PlantsRelation PlantFolderRelation
}

type PlantFolderRelation struct {
	gorm.Model
	FolderID uint64
	PlantsID pq.Int64Array `gorm:"type:integer[]"`
}
