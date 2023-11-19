package gormModels

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

func (pf *PlantFolderRelation) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&Folder{}).
		Where("id = ?", pf.FolderID).
		UpdateColumn("plants_num", gorm.Expr("plants_num + ?", 1))
	return
}

func (pf *PlantFolderRelation) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&Folder{}).
		Where("id = ?", pf.FolderID).
		UpdateColumn("plants_num", gorm.Expr("?", uint64(len(pf.PlantsID))))
	return
}
