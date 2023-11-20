package domain

import (
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type PlantsUsecase interface {
	AddPlant(plant httpModels.Plant) error
	GetPlant(userID uint64, plantID int64) (httpModels.Plant, error)
	GetPlants(userID uint64) ([]httpModels.XiaomiPlant, error)
	DeletePlant(userID, plantID uint64) error
	GetPlantByID(plantID uint64) (gormModels.XiaomiPlant, error)
	GetPlantByName(plantName string) ([]httpModels.XiaomiPlant, error)
	GetPlantsPage(pageNum uint64) ([]httpModels.XiaomiPlant, error)

	CreateCustomPlant(plant httpModels.CustomPlant, extension string) (uint64, error)
	UpdateCustomPlant(plant httpModels.CustomPlant, extension string) error
	GetCustomPlants(userID uint64) ([]httpModels.CustomPlant, error)
	GetCustomPlant(userID, plantID uint64) (httpModels.CustomPlant, error)
	DeleteCustomPlant(userID, plantID uint64) error
	GetCustomPlantImage(userID, plantID uint64) (string, error)
}
