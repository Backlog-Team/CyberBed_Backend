package gormModels

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint64 `json:"userID"   gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Creating default folders after creating user
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
  tx.Create(&DefaultFolder{
    UserID: u.ID,
    FolderName: DefaultFolderName,
  })
  return
}
