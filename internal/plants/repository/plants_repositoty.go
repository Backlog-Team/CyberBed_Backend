package plantsRepository

import (
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"

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
		&models.XiaomiPlant{},
		&models.XiaomiPlantBasic{},
		&models.XiaomiPlantMaintenance{},
		&models.XiaomiPlantPrameter{},
		&models.UserPlants{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateUserPlantsRelations(userID uint64, plantsID []int64) error {
	res := db.DB.Create(&models.UserPlants{
		UserID:   userID,
		PlantsID: pq.Int64Array(plantsID),
	})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (db *Postgres) AddUserPlantsRelations(userID uint64, plantsID []int64) error {
	userPlant := []models.UserPlants{}
	db.DB.Table(models.PlantsTable).Select("*").Where("user_id = ?", userID).Scan(&userPlant)

	if len(userPlant) == 0 {
		res := db.DB.Table(models.PlantsTable).Create(&models.UserPlants{
			UserID:   userID,
			PlantsID: pq.Int64Array(plantsID),
		})
		if res.Error != nil {
			return res.Error
		}
	} else {
		newPlantIDs := userPlant[0].PlantsID
		newPlantIDs = append(newPlantIDs, plantsID...)

		res := db.DB.Table(models.PlantsTable).Where("user_id = ?", userID).Update("plants_id", &newPlantIDs)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (db *Postgres) GetPlantsByUserID(userID uint64) (models.UserPlants, error) {
	var pl models.UserPlants
	if err := db.DB.Table(models.PlantsTable).
		Select("*").
		Where("user_id = ?", userID).
		Scan(&pl).
		Error; err != nil {
		return models.UserPlants{}, err
	}

	return pl, nil
}

func (db *Postgres) UpdateUserPlantsRelation(relation models.UserPlants) error {
	if err := db.DB.Table(models.PlantsTable).
		Where("user_id = ?", relation.UserID).
		Update("plants_id", &relation.PlantsID).Error; err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetPlantByID(plantID uint64) (models.XiaomiPlant, error) {
	var plant models.XiaomiPlant
	if err := db.DB.Preload("Basic").
		Preload("Maintenance").
		Preload("Parameter").
		Where("xiaomi_plants.id = ?", plantID).
		First(&plant).Error; err != nil {
		return models.XiaomiPlant{}, err
	}
	return plant, nil
}

func (db *Postgres) GetByPlantName(plantName string) ([]models.XiaomiPlant, error) {
	var plants []models.XiaomiPlant
	if err := db.DB.Preload("Basic").
		Preload("Maintenance").
		Preload("Parameter").
		Where("pid LIKE ?", "%"+strings.ToLower(plantName)+"%").
		Limit(10).
		Find(&plants).Error; err != nil {
		return nil, err
	}
	return plants, nil
}

func (db *Postgres) GetPlantsPage(pageNum uint64) ([]models.XiaomiPlant, error) {
	pageSize := 10
	offset := (pageSize - 1) * int(pageNum-1)
	var plants []models.XiaomiPlant
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
