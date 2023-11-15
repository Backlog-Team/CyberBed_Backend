package foldersRepository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/models"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&models.Folder{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateFolder(folder models.Folder) (uint64, error) {
	var res models.Folder
	if err := db.DB.Model(&models.Folder{}).
		Create(folder).
		Scan(&res).
		Error; err != nil {
		return 0, err
	}
	return res.ID, nil
}

func (db *Postgres) GetFolders(userID uint64) ([]models.Folder, error) {
	var resRows []models.Folder
	if err := db.DB.Preload("folders").
		Where("user_id = ?", userID).
		Find(&resRows).
		Error; err != nil {
		return nil, err
	}
	return resRows, nil
}

func (db *Postgres) GetFolder(id uint64) (models.Folder, error) {
	var folderRow models.Folder
	if err := db.DB.Model(&models.Folder{}).
		Where("id = ?", id).
		First(&folderRow).
		Error; err != nil {
		return models.Folder{}, err
	}
	return folderRow, nil
}

func (db *Postgres) DeleteFolder(id uint64) error {
	if err := db.DB.Model(&models.Folder{}).
		Where("id = ?", id).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetPlantsID(folderID uint64) ([]uint64, error) {
	var plantIDs []int64
	if err := db.DB.Table("plant_folder_relations").
		Select("plants_id").
		Where("folder_id = ?", folderID).
		First(&plantIDs).
		Error; err != nil {
		return nil, err
	}

	var convertedPlants []uint64
	for i, plant := range convertedPlants {
		convertedPlants[i] = uint64(plant)
	}

	return convertedPlants, nil
}
