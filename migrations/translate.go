package migrations

import (
	"encoding/json"
	"log"
	"os"

	gormModels "github.com/cyber_bed/internal/models/gorm"
	"github.com/cyber_bed/internal/utils/translator"
)

const (
	pathToDir = "/home/milchenko/technopark/fourth_semestr/CyberBed_Backend/migrations/plant-database/json_rus/"
)

func EditField() error {
	entries, err := os.ReadDir(pathToDir)
	if err != nil {
		return err
	}

	for plantIndx, e := range entries {
		content, err := os.ReadFile(pathToDir + e.Name())
		if err != nil {
			return err
		}
		var plantItem gormModels.XiaomiPlant
		json.Unmarshal(content, &plantItem)

		plantItem.Basic.FloralLanguage, err = translator.Translate(plantItem.Basic.FloralLanguage)
		if err != nil {
			return err
		}

		contentToWrite, err := json.Marshal(plantItem)
		if err != nil {
			return err
		}

		if err = os.WriteFile(pathToDir+e.Name(), contentToWrite, 0644); err != nil {
			return err
		}

		log.Printf("Edited %d files", plantIndx+1)
	}

	return nil
}

func TranslatePlants(pathToDir string) error {
	entries, err := os.ReadDir(pathToDir)
	filesNum := len(entries)
	if err != nil {
		return err
	}

	if _, err = os.Stat(pathToDir + "/../json_rus"); err != nil {
		if err := os.Mkdir(pathToDir+"/../json_rus", 0755); err != nil {
			return err
		}
	}

	for plantIndx, e := range entries {
		content, err := os.ReadFile(pathToDir + "/" + e.Name())
		if err != nil {
			return err
		}

		var plantItem gormModels.XiaomiPlant
		json.Unmarshal(content, &plantItem)

		plantItem.DisplayPid, err = translator.Translate(plantItem.DisplayPid)
		if err != nil {
			return err
		}
		plantItem.Basic.Origin, err = translator.Translate(plantItem.Basic.Origin)
		if err != nil {
			return err
		}
		plantItem.Basic.Production, err = translator.Translate(plantItem.Basic.Production)
		if err != nil {
			return err
		}
		plantItem.Basic.Category, err = translator.Translate(plantItem.Basic.Category)
		if err != nil {
			return err
		}
		plantItem.Basic.Blooming, err = translator.Translate(plantItem.Basic.Blooming)
		if err != nil {
			return err
		}
		plantItem.Basic.Color, err = translator.Translate(plantItem.Basic.Color)
		if err != nil {
			return err
		}

		plantItem.Maintenance.Size, err = translator.Translate(plantItem.Maintenance.Size)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Soil, err = translator.Translate(plantItem.Maintenance.Soil)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Sunlight, err = translator.Translate(plantItem.Maintenance.Sunlight)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Watering, err = translator.Translate(plantItem.Maintenance.Watering)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Fertilization, err = translator.Translate(
			plantItem.Maintenance.Fertilization,
		)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Pruning, err = translator.Translate(plantItem.Maintenance.Pruning)
		if err != nil {
			return err
		}

		textToWrite, err := json.Marshal(plantItem)
		if err != nil {
			return err
		}

		if err = os.WriteFile(pathToDir+"/../json_rus/"+plantItem.PlantID+".json", textToWrite, 0644); err != nil {
			return err
		}

		log.Printf("written %d file of %d", plantIndx+1, filesNum)
	}

	return nil
}
