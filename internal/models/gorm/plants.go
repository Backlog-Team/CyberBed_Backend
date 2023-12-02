package gormModels

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserPlants struct {
	UserID   uint64
	PlantsID pq.Int64Array `gorm:"type:integer[]"`
}

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
	gorm.Model
	ID        uint64
	UserID    uint64
	PlantName string
	About     string
	Image     string
	Notify    time.Time
}

type SavedPlant struct {
	gorm.Model
	ID      uint64
	UserID  uint64
	PlantID uint64
}

func (sp *SavedPlant) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("id = ?", sp.ID).Delete(&Channel{})
	tx.Where("id = ?", sp.ID).Delete(&Notification{})
	return
}

type Channel struct {
	gorm.Model
	UserID    uint64
	PlantID   uint64
	ChannelID uint64
}

type PlantStat struct {
	gorm.Model
	PlantID uint64
	UserID  uint64
	IsLiked bool `gorm:"default:false"`
	IsSaved bool `gorm:"default:false"`
}

type XiaomiPlant struct {
	gorm.Model
	ID          uint64
	PlantID     string                 `json:"pid"`
	Basic       XiaomiPlantBasic       `json:"basic"       gorm:"foreignkey:XiaomiPlantID;association_foreignkey:ID"`
	DisplayPid  string                 `json:"display_pid" gorm:"index"`
	Maintenance XiaomiPlantMaintenance `json:"maintenance" gorm:"foreignkey:XiaomiPlantID;association_foreignkey:ID"`
	Parameter   XiaomiPlantPrameter    `json:"parameter"   gorm:"foreignkey:XiaomiPlantID;association_foreignkey:ID"`
	Image       string                 `json:"image"`
}

type XiaomiPlantBasic struct {
	gorm.Model
	XiaomiPlantID  uint64 `gorm:"index"`
	FloralLanguage string `             json:"floral_language"`
	Origin         string `             json:"origin"`
	Production     string `             json:"production"`
	Category       string `             json:"category"`
	Blooming       string `             json:"blooming"`
	Color          string `             json:"color"`
}

type XiaomiPlantMaintenance struct {
	gorm.Model
	XiaomiPlantID uint64 `gorm:"index"`
	Size          string `             json:"size"`
	Soil          string `             json:"soil"`
	Sunlight      string `             json:"sunlight"`
	Watering      string `             json:"watering"`
	Fertilization string `             json:"fertilization"`
	Pruning       string `             json:"pruning"`
}

type XiaomiPlantPrameter struct {
	gorm.Model
	XiaomiPlantID uint64 `gorm:"index"`
	MaxLightMmol  uint64 `             json:"max_light_mmol"`
	MinLightMmol  uint64 `             json:"min_light_mmol"`
	MaxLightLux   uint64 `             json:"max_light_lux"`
	MinLightLux   uint64 `             json:"min_light_lux"`
	MaxTemp       uint64 `             json:"max_temp"`
	MinTemp       uint64 `             json:"min_temp"`
	MaxEnvHumid   uint64 `             json:"max_env_humidity"`
	MinEnvHumid   uint64 `             json:"min_env_humidity"`
	MaxSoilMoist  uint64 `             json:"max_soil_moisture"`
	MinSoilMoist  uint64 `             json:"min_soil_moisture"`
	MaxSoilEc     uint64 `             json:"max_soil_ec"`
	MinSoilEc     uint64 `             json:"min_soil_ec"`
}
