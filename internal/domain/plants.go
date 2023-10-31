package domain

import "github.com/cyber_bed/internal/models"

type PlantsUsecase interface {
	AddPlant(plant models.Plant) error
	GetPlant(userID uint64, plantID int64) (models.Plant, error)
	GetPlants(userID uint64) ([]models.XiaomiPlant, error)
	DeletePlant(userID, plantID uint64) error
	GetPlantByID(plantID uint64) (models.XiaomiPlant, error)
	GetPlantByName(plantName string) ([]models.XiaomiPlant, error)
	GetPlantsPage(pageNum uint64) ([]models.XiaomiPlant, error)
}
