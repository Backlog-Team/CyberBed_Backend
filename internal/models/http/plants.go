package httpModels

import (
	gormModels "github.com/cyber_bed/internal/models/gorm"
)

type Plant struct {
	UserID   uint64 `json:"userID"`
	ID       uint64 `json:"id"`
	ImageUrl string `json:"imageUrl"`

	CommonName     string   `json:"common_name"`
	ScientificName []string `json:"scientific_name"`
	OtherName      []string `json:"other_name"`
	Cycle          string   `json:"cycle"`
	Watering       string   `json:"watering"`
}

type CustomPlant struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"userID"`
	PlantName    string `json:"plantName"`
	About        string `json:"about"`
	Image        string `json:"image"`
	NotifyString string `json:"notify"`
}

type SavedPlant struct {
	UserID  uint64
	PlantID uint64
}

type XiaomiPlant struct {
	ID      uint64
	PlantID string `json:"pid"`
	Basic   struct {
		FloralLanguage string `             json:"floral_language"`
		Origin         string `             json:"origin"`
		Production     string `             json:"production"`
		Category       string `             json:"category"`
		Blooming       string `             json:"blooming"`
		Color          string `             json:"color"`
	} `json:"basic"`
	DisplayPid  string `json:"display_pid"`
	Maintenance struct {
		Size          string `             json:"size"`
		Soil          string `             json:"soil"`
		Sunlight      string `             json:"sunlight"`
		Watering      string `             json:"watering"`
		Fertilization string `             json:"fertilization"`
		Pruning       string `             json:"pruning"`
	} `json:"maintenance"`
	Parameter struct {
		MaxLightMmol uint64 `             json:"max_light_mmol"`
		MinLightMmol uint64 `             json:"min_light_mmol"`
		MaxLightLux  uint64 `             json:"max_light_lux"`
		MinLightLux  uint64 `             json:"min_light_lux"`
		MaxTemp      uint64 `             json:"max_temp"`
		MinTemp      uint64 `             json:"min_temp"`
		MaxEnvHumid  uint64 `             json:"max_env_humidity"`
		MinEnvHumid  uint64 `             json:"min_env_humidity"`
		MaxSoilMoist uint64 `             json:"max_soil_moisture"`
		MinSoilMoist uint64 `             json:"min_soil_moisture"`
		MaxSoilEc    uint64 `             json:"max_soil_ec"`
		MinSoilEc    uint64 `             json:"min_soil_ec"`
	} `json:"parameter"`
	IsLiked     bool     `json:"isLiked"`
	IsSaved     bool     `json:"isSaved"`
	FolderSaved []Folder `json:"folderSaved"`
}

func XiaomiPlantGormToHttp(plant gormModels.XiaomiPlant) XiaomiPlant {
	res := XiaomiPlant{
		ID:         plant.ID,
		PlantID:    plant.PlantID,
		DisplayPid: plant.DisplayPid,
		Basic: struct {
			FloralLanguage string `             json:"floral_language"`
			Origin         string `             json:"origin"`
			Production     string `             json:"production"`
			Category       string `             json:"category"`
			Blooming       string `             json:"blooming"`
			Color          string `             json:"color"`
		}{
			FloralLanguage: plant.Basic.FloralLanguage,
			Origin:         plant.Basic.Origin,
			Production:     plant.Basic.Production,
			Category:       plant.Basic.Category,
			Blooming:       plant.Basic.Blooming,
			Color:          plant.Basic.Color,
		},
		Maintenance: struct {
			Size          string `             json:"size"`
			Soil          string `             json:"soil"`
			Sunlight      string `             json:"sunlight"`
			Watering      string `             json:"watering"`
			Fertilization string `             json:"fertilization"`
			Pruning       string `             json:"pruning"`
		}{
			Size:          plant.Maintenance.Size,
			Soil:          plant.Maintenance.Soil,
			Sunlight:      plant.Maintenance.Sunlight,
			Watering:      plant.Maintenance.Watering,
			Fertilization: plant.Maintenance.Fertilization,
			Pruning:       plant.Maintenance.Pruning,
		},
		Parameter: struct {
			MaxLightMmol uint64 `             json:"max_light_mmol"`
			MinLightMmol uint64 `             json:"min_light_mmol"`
			MaxLightLux  uint64 `             json:"max_light_lux"`
			MinLightLux  uint64 `             json:"min_light_lux"`
			MaxTemp      uint64 `             json:"max_temp"`
			MinTemp      uint64 `             json:"min_temp"`
			MaxEnvHumid  uint64 `             json:"max_env_humidity"`
			MinEnvHumid  uint64 `             json:"min_env_humidity"`
			MaxSoilMoist uint64 `             json:"max_soil_moisture"`
			MinSoilMoist uint64 `             json:"min_soil_moisture"`
			MaxSoilEc    uint64 `             json:"max_soil_ec"`
			MinSoilEc    uint64 `             json:"min_soil_ec"`
		}{
			MaxLightMmol: plant.Parameter.MaxLightMmol,
			MinLightMmol: plant.Parameter.MinLightMmol,
			MaxLightLux:  plant.Parameter.MaxLightLux,
			MinLightLux:  plant.Parameter.MinLightLux,
			MaxTemp:      plant.Parameter.MaxTemp,
			MinTemp:      plant.Parameter.MinTemp,
			MaxEnvHumid:  plant.Parameter.MaxEnvHumid,
			MinEnvHumid:  plant.Parameter.MinEnvHumid,
			MaxSoilMoist: plant.Parameter.MaxSoilMoist,
			MinSoilMoist: plant.Parameter.MinSoilMoist,
			MaxSoilEc:    plant.Parameter.MaxSoilEc,
			MinSoilEc:    plant.Parameter.MinSoilEc,
		},
	}
	return res
}

func CustomPlantGormToHttp(plant gormModels.CustomPlant) CustomPlant {
	return CustomPlant{
		ID:        plant.ID,
		UserID:    plant.UserID,
		PlantName: plant.PlantName,
		About:     plant.About,
	}
}

func SavedPlantGormToHttp(plant gormModels.SavedPlant) SavedPlant {
	return SavedPlant{
		UserID:  plant.UserID,
		PlantID: plant.PlantID,
	}
}
