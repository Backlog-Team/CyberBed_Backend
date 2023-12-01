package foldersRepository

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
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
		&gormModels.Folder{},
		&gormModels.PlantFolderRelation{},

		&gormModels.Channel{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateFolder(folder httpModels.Folder) (uint64, error) {
	var res gormModels.Folder
	if err := db.DB.Model(&gormModels.Folder{}).
		Create(&gormModels.Folder{
			FolderName: folder.FolderName,
			UserID:     folder.UserID,
			IsDefalut:  false,
		}).
		Scan(&res).
		Error; err != nil {
		return 0, err
	}
	return res.ID, nil
}

func (db *Postgres) GetFolders(userID uint64) ([]gormModels.Folder, error) {
	var resRows []gormModels.Folder
	if err := db.DB.Model(&gormModels.Folder{}).
		Where("user_id = ?", userID).
		Find(&resRows).
		Error; err != nil {
		return nil, err
	}
	return resRows, nil
}

func (db *Postgres) GetFolder(id uint64) (gormModels.Folder, error) {
	var folderRow gormModels.Folder
	if err := db.DB.Model(&gormModels.Folder{}).
		Where("id = ?", id).
		First(&folderRow).
		Error; err != nil {
		return gormModels.Folder{}, err
	}
	return folderRow, nil
}

func (db *Postgres) GetFolderByNameAndUserID(
	folderName string,
	userID uint64,
) (gormModels.Folder, error) {
	var folderRow gormModels.Folder
	if err := db.DB.Model(&gormModels.Folder{}).
		Where("folder_name = ? AND user_id = ?", folderName, userID).
		First(&folderRow).
		Error; err != nil {
		return gormModels.Folder{}, err
	}
	return folderRow, nil
}

func (db *Postgres) DeleteFolder(id uint64) error {
	if err := db.DB.Select("Folder").
		Where("id = ? AND is_default = false", id).
		Delete(&gormModels.Folder{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetPlantsID(folderID uint64) ([]uint64, error) {
	var plantIDs gormModels.PlantFolderRelation
	if err := db.DB.Model(&gormModels.PlantFolderRelation{}).
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
	folderPlant := gormModels.PlantFolderRelation{}
	if err := db.DB.Model(&gormModels.PlantFolderRelation{}).
		Where("folder_id = ?", folderID).
		First(&folderPlant).
		Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			if err := db.DB.Model(&gormModels.PlantFolderRelation{}).
				Create(&gormModels.PlantFolderRelation{
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
	if err := db.DB.Model(&gormModels.PlantFolderRelation{FolderID: folderID}).
		Where("folder_id = ?", folderID).
		Update("plants_id", &folderPlant.PlantsID).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) UpdateFolderPlant(folderID, plantID uint64) error {
	folderPlant := gormModels.PlantFolderRelation{}
	if err := db.DB.Model(&gormModels.PlantFolderRelation{}).
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

	if err := db.DB.Model(&gormModels.PlantFolderRelation{FolderID: folderID}).
		Where("folder_id = ?", folderID).
		Update("plants_id", &folderPlant.PlantsID).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetFolderByPlantAndUserID(userID, plantID uint64) ([]gormModels.Folder, error) {
	var folderRow []gormModels.Folder
	if err := db.DB.Model(&gormModels.PlantFolderRelation{}).
		Joins("JOIN folders ON plant_folder_relations.folder_id=folders.id").
		Where("user_id = ? AND plants_id @> ARRAY[?]::integer[]", userID, plantID).
		Find(&folderRow).Error; err != nil {
		return []gormModels.Folder{}, err
	}
	return folderRow, nil
}

func (db *Postgres) CreateChannel(folderID, plantID, channelID uint64) (uint64, error) {
	var chanRow gormModels.Channel
	if err := db.DB.Create(&gormModels.Channel{
		FolderID:  folderID,
		PlantID:   plantID,
		ChannelID: channelID,
	}).Scan(&chanRow).Error; err != nil {
		return 0, err
	}
	return uint64(chanRow.ID), nil
}

func (db *Postgres) GetChannelByFolderPlantID(folderID, plantID uint64) (uint64, error) {
	var chanRow gormModels.Channel
	if err := db.DB.Model(&gormModels.Channel{}).
		Where("folder_id = ? AND plant_id = ?", folderID, plantID).
		First(&chanRow).Error; err != nil {
		return 0, err
	}
	return chanRow.ChannelID, nil
}
