package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	ID             uint64 `gorm:"primaryKey"`
	UserID         uint64
	FolderName     string
	PlantsNum      uint64
	PlantsRelation PlantFolderRelation
}

type PlantFolderRelation struct {
	gorm.Model
	FolderID uint64        `gorm:"index;unique"`
	PlantsID pq.Int64Array `gorm:"type:integer[]"`
}

type FolderHttp struct {
	ID         uint64 `json:"ID"`
	UserID     uint64 `json:"userID"`
	FolderName string `json:"folderName"`
	PlantsNum  uint64 `json:"plantsNum"`
}

func FolderGormToHttp(f Folder) FolderHttp {
	return FolderHttp{
		ID:         f.ID,
		UserID:     f.UserID,
		FolderName: f.FolderName,
		PlantsNum:  f.PlantsNum,
	}
}

func (pf *PlantFolderRelation) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&Folder{}).
		Where("id = ?", pf.FolderID).
		UpdateColumn("plants_num", gorm.Expr("plants_num + ?", 1))
	return
}

// FIXME: this hook doesn't work right way
func (pf *PlantFolderRelation) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&Folder{}).
		Where("id = ?", pf.FolderID).
		UpdateColumn("plants_num", uint64(len(pf.PlantsID)))
	return
}
