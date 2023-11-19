package plantsRepository

import (
	"strings"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormModels "github.com/cyber_bed/internal/models/gorm"
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
		newPlantIDs := userPlant[0].PlantsID
		newPlantIDs = append(newPlantIDs, plantsID...)

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
		Where("plant_id LIKE ? OR display_pid LIKE ?", "%"+strings.ToLower(plantName)+"%", "%"+strings.ToLower(plantName)+"%").
		Limit(10).
		Find(&plants).Error; err != nil {
		return nil, err
	}
	return plants, nil
}

func (db *Postgres) GetPlantsPage(pageNum uint64) ([]gormModels.XiaomiPlant, error) {
	pageSize := 10
	offset := (pageSize - 1) * int(pageNum-1)
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
