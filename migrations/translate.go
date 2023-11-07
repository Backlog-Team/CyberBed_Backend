package migrations

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	translator "github.com/turk/free-google-translate"

	"github.com/cyber_bed/internal/models"
)

func Translate(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	client := http.Client{}
	t := translator.NewTranslator(&client)
	result, err := t.Translate(text, "en", "ru")
	if err != nil {
		return "", err
	}

	return result, nil
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

		var plantItem models.XiaomiPlant
		json.Unmarshal(content, &plantItem)

		plantItem.DisplayPid, err = Translate(plantItem.DisplayPid)
		if err != nil {
			return err
		}
		plantItem.Basic.Origin, err = Translate(plantItem.Basic.Origin)
		if err != nil {
			return err
		}
		plantItem.Basic.Production, err = Translate(plantItem.Basic.Production)
		if err != nil {
			return err
		}
		plantItem.Basic.Category, err = Translate(plantItem.Basic.Category)
		if err != nil {
			return err
		}
		plantItem.Basic.Blooming, err = Translate(plantItem.Basic.Blooming)
		if err != nil {
			return err
		}
		plantItem.Basic.Color, err = Translate(plantItem.Basic.Color)
		if err != nil {
			return err
		}

		plantItem.Maintenance.Size, err = Translate(plantItem.Maintenance.Size)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Soil, err = Translate(plantItem.Maintenance.Soil)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Sunlight, err = Translate(plantItem.Maintenance.Sunlight)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Watering, err = Translate(plantItem.Maintenance.Watering)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Fertilization, err = Translate(plantItem.Maintenance.Fertilization)
		if err != nil {
			return err
		}
		plantItem.Maintenance.Pruning, err = Translate(plantItem.Maintenance.Pruning)
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
