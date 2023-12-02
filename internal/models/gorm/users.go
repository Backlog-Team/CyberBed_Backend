package gormModels

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint64 `json:"userID"   gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Creating default folder after creating user
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	tx.Create(&Folder{
		UserID:         u.ID,
		FolderName:     DefaultFolderName,
		IsDefalut:      true,
		PlantsRelation: PlantFolderRelation{},
	})
	return
}
