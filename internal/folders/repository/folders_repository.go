package foldersRepository

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
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
		&models.PlantFolderRelation{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateFolder(folder models.Folder) (uint64, error) {
	var res models.Folder
	if err := db.DB.Model(&models.Folder{}).
		Create(&folder).
		Scan(&res).
		Error; err != nil {
		return 0, err
	}
	return res.ID, nil
}

func (db *Postgres) GetFolders(userID uint64) ([]models.Folder, error) {
	var resRows []models.Folder
	if err := db.DB.Model(&models.Folder{}).
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
	if err := db.DB.Select("Folder").
		Delete(&models.Folder{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetPlantsID(folderID uint64) ([]uint64, error) {
	var plantIDs models.PlantFolderRelation
	if err := db.DB.Model(&models.PlantFolderRelation{}).
		Where("folder_id = ?", folderID).
		First(&plantIDs).
		Error; err != nil {
		return nil, err
	}

	var convertedPlants []uint64
	for _, plant := range plantIDs.PlantsID {
		convertedPlants = append(convertedPlants, uint64(plant))
	}

	return convertedPlants, nil
}

func (db *Postgres) AddPlantToFolder(folderID, plantID uint64) error {
	folderPlant := models.PlantFolderRelation{}
	if err := db.DB.Model(&models.PlantFolderRelation{}).
		Where("folder_id = ?", folderID).
		First(&folderPlant).
		Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			if err := db.DB.Model(&models.PlantFolderRelation{}).
				Create(&models.PlantFolderRelation{
					FolderID: folderID,
					PlantsID: pq.Int64Array{int64(plantID)},
				}).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	folderPlant.PlantsID = append(folderPlant.PlantsID, int64(plantID))
	if err := db.DB.Model(&models.PlantFolderRelation{}).
		Where("folder_id = ?", folderID).
		Update("plants_id", &folderPlant.PlantsID).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) UpdateFolderPlant(folderID, plantID uint64) error {
	folderPlant := models.PlantFolderRelation{}
	if err := db.DB.Model(&models.PlantFolderRelation{}).
		Where("folder_id = ? AND plants_id @> ARRAY[?]::integer[]", folderID, plantID).
		First(&folderPlant).
		Error; err != nil {
		return err
	}

	for index, pID := range folderPlant.PlantsID {
		if uint64(pID) == plantID {
			folderPlant.PlantsID = append(
				folderPlant.PlantsID[:index],
				folderPlant.PlantsID[index+1:]...,
			)
			break
		}
	}

	if err := db.DB.Model(&models.PlantFolderRelation{}).
		Where("folder_id = ?", folderID).
		Update("plants_id", &folderPlant.PlantsID).
		Error; err != nil {
		return err
	}
	return nil
}
