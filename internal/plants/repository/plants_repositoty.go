package plantsRepository

import (
	"github.com/lib/pq"
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
		&gormModels.XiaomiPlant{},
		&gormModels.XiaomiPlantBasic{},
		&gormModels.XiaomiPlantMaintenance{},
		&gormModels.XiaomiPlantPrameter{},
		&gormModels.UserPlants{},

		&gormModels.CustomPlant{},
		&gormModels.SavedPlant{},
		&gormModels.PlantStat{},

		&gormModels.Channel{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateUserPlantsRelations(userID uint64, plantsID []int64) error {
	res := db.DB.Create(&gormModels.UserPlants{
		UserID:   userID,
		PlantsID: pq.Int64Array(plantsID),
	})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (db *Postgres) AddUserPlantsRelations(userID uint64, plantsID []int64) error {
	userPlant := []gormModels.UserPlants{}
	db.DB.Table(gormModels.PlantsTable).Select("*").Where("user_id = ?", userID).Scan(&userPlant)

	if len(userPlant) == 0 {
		res := db.DB.Table(gormModels.PlantsTable).Create(&gormModels.UserPlants{
			UserID:   userID,
			PlantsID: pq.Int64Array(plantsID),
		})
		if res.Error != nil {
			return res.Error
		}
	} else {
		plantsMap := make(map[uint64]bool)
		for _, pid := range userPlant[0].PlantsID {
			plantsMap[uint64(pid)] = true
		}
		for _, pid := range plantsID {
			plantsMap[uint64(pid)] = true
		}

		newPlantIDs := make(pq.Int64Array, 0)
		for key := range plantsMap {
			newPlantIDs = append(newPlantIDs, int64(key))
		}

		res := db.DB.Table(gormModels.PlantsTable).Where("user_id = ?", userID).Update("plants_id", &newPlantIDs)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (db *Postgres) GetPlantsByUserID(userID uint64) (gormModels.UserPlants, error) {
	var pl gormModels.UserPlants
	if err := db.DB.Table(gormModels.PlantsTable).
		Select("*").
		Where("user_id = ?", userID).
		Scan(&pl).
		Error; err != nil {
		return gormModels.UserPlants{}, err
	}

	return pl, nil
}

func (db *Postgres) UpdateUserPlantsRelation(relation gormModels.UserPlants) error {
	if err := db.DB.Table(gormModels.PlantsTable).
		Where("user_id = ?", relation.UserID).
		Update("plants_id", &relation.PlantsID).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetPlantByID(plantID uint64) (gormModels.XiaomiPlant, error) {
	var plant gormModels.XiaomiPlant
	if err := db.DB.Preload("Basic").
		Preload("Maintenance").
		Preload("Parameter").
		Where("xiaomi_plants.id = ?", plantID).
		First(&plant).Error; err != nil {
		return gormModels.XiaomiPlant{}, err
	}
	return plant, nil
}

func (db *Postgres) GetByPlantName(plantName string) ([]gormModels.XiaomiPlant, error) {
	var plants []gormModels.XiaomiPlant
	if err := db.DB.Preload("Basic").
		Preload("Maintenance").
		Preload("Parameter").
		Where("plant_id LIKE LOWER(?) OR LOWER(display_pid) LIKE LOWER(?) ", "%"+plantName+"%", "%"+plantName+"%").
		Limit(10).
		Find(&plants).Error; err != nil {
		return nil, err
	}
	return plants, nil
}

func (db *Postgres) GetPlantsPage(pageNum uint64) ([]gormModels.XiaomiPlant, error) {
	pageSize := 10
	offset := pageSize * int(pageNum-1)
	var plants []gormModels.XiaomiPlant
	if err := db.DB.Preload("Basic").
		Preload("Maintenance").
		Preload("Parameter").
		Order("id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&plants).
		Error; err != nil {
		return nil, err
	}
	return plants, nil
}

func (db *Postgres) CreateCustomPlant(plant httpModels.CustomPlant) (uint64, error) {
	var plantRow gormModels.CustomPlant
	if err := db.DB.Create(&gormModels.CustomPlant{
		PlantName: plant.PlantName,
		About:     plant.About,
		UserID:    plant.UserID,
		Image:     plant.Image,
	}).Scan(&plantRow).Error; err != nil {
		return 0, err
	}
	return plantRow.ID, nil
}

func (db *Postgres) UpdateCustomPlant(plant httpModels.CustomPlant) error {
	if err := db.DB.Model(&gormModels.CustomPlant{}).
		Where("id = ?", plant.ID).
		Update("plant_name", plant.PlantName).
		Update("about", plant.About).
		Update("image", plant.Image).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetCustomPlants(userID uint64) ([]gormModels.CustomPlant, error) {
	var resRows []gormModels.CustomPlant
	if err := db.DB.Model(&gormModels.CustomPlant{}).
		Where("user_id = ?", userID).
		Find(&resRows).Error; err != nil {
		return nil, err
	}
	return resRows, nil
}

func (db *Postgres) GetCustomPlant(userID, plantID uint64) (gormModels.CustomPlant, error) {
	var resRow gormModels.CustomPlant
	if err := db.DB.Model(&gormModels.CustomPlant{}).
		Where("user_id = ? AND id = ?", userID, plantID).
		First(&resRow).Error; err != nil {
		return gormModels.CustomPlant{}, err
	}
	return resRow, nil
}

func (db *Postgres) DeleteCustomPlant(userID, plantID uint64) error {
	if err := db.DB.Select("CustomPlant").
		Delete(&gormModels.CustomPlant{
			ID:     plantID,
			UserID: userID,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) CreateSavedPlant(userID, plantID uint64) error {
	if err := db.DB.Create(&gormModels.SavedPlant{
		UserID:  userID,
		PlantID: plantID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) DeleteSavedPlant(userID, plantID uint64) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Select("SavedPlant").
		Where("user_id = ? AND plant_id = ?", userID, plantID).
		Delete(&gormModels.SavedPlant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ? AND plant_id = ?", userID, plantID).
		Delete(&gormModels.Channel{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ? AND plant_id = ?", userID, plantID).
		Delete(&gormModels.Notification{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *Postgres) GetSavedPlants(userID uint64) ([]gormModels.SavedPlant, error) {
	var plantRows []gormModels.SavedPlant
	if err := db.DB.Model(&gormModels.SavedPlant{
		UserID: userID,
	}).Find(&plantRows).Error; err != nil {
		return []gormModels.SavedPlant{}, err
	}
	return plantRows, nil
}

func (db *Postgres) GetSavedPlantByIDs(userID, plantID uint64) (gormModels.SavedPlant, error) {
	var plantRow gormModels.SavedPlant
	if err := db.DB.Model(&gormModels.SavedPlant{}).
		Where("user_id = ? AND plant_id = ?", userID, plantID).
		First(&plantRow).Error; err != nil {
		return gormModels.SavedPlant{}, err
	}
	return plantRow, nil
}

func (db *Postgres) CreateChannel(plantID, channelID, userID uint64) (uint64, error) {
	var chanRow gormModels.Channel
	if err := db.DB.Create(&gormModels.Channel{
		UserID:    userID,
		PlantID:   plantID,
		ChannelID: channelID,
	}).Scan(&chanRow).Error; err != nil {
		return 0, err
	}
	return uint64(chanRow.ID), nil
}

func (db *Postgres) GetChannelByUserAndPlantID(userID, plantID uint64) (uint64, error) {
	var chanRow gormModels.Channel
	if err := db.DB.Model(&gormModels.Channel{}).
		Where("user_id = ? AND plant_id = ?", userID, plantID).
		First(&chanRow).Error; err != nil {
		return 0, err
	}
	return chanRow.ChannelID, nil
}

func (db *Postgres) UpdateChannelByUserAndPlantID(userID, plantID, channelID uint64) error {
	return db.DB.Model(&gormModels.Channel{}).
		Where("user_id = ? AND plant_id = ?", userID, plantID).
		Update("channel_id", channelID).
		Error
}
