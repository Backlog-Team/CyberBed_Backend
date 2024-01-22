package migrations

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	gormModels "github.com/cyber_bed/internal/models/gorm"
	"github.com/cyber_bed/migrations/translator"
)

//const (
//	pathToDir = "/home/milchenko/technopark/fourth_semestr/CyberBed_Backend/migrations/plant-database/json_rus/"
//)
//
//func EditField() error {
//	entries, err := os.ReadDir(pathToDir)
//	if err != nil {
//		return err
//	}
//
//	for plantIndx, e := range entries {
//		content, err := os.ReadFile(pathToDir + e.Name())
//		if err != nil {
//			return err
//		}
//		var plantItem gormModels.XiaomiPlant
//		json.Unmarshal(content, &plantItem)
//
//		plantItem.Basic.FloralLanguage, err = translator.Translate(plantItem.Basic.FloralLanguage)
//		if err != nil {
//			return err
//		}
//
//		contentToWrite, err := json.Marshal(plantItem)
//		if err != nil {
//			return err
//		}
//
//		if err = os.WriteFile(pathToDir+e.Name(), contentToWrite, 0644); err != nil {
//			return err
//		}
//
//		log.Printf("Edited %d files", plantIndx+1)
//	}
//
//	return nil
//}

func TranslatePlants(pathToDir, apiKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	entries, err := os.ReadDir(pathToDir)
	filesNum := len(entries)
	if err != nil {
		return err
	}

	if _, err = os.Stat(pathToDir + "/../json_rus"); err != nil {
		if err = os.Mkdir(pathToDir+"/../json_rus", 0755); err != nil {
			return err
		}
	}

	eg, ctx := errgroup.WithContext(ctx)

	for plantIndx, e := range entries {
		index := plantIndx
		entry := e

		eg.Go(func() error {
			content, err := os.ReadFile(pathToDir + "/" + entry.Name())
			if err != nil {
				return err
			}

			var plantItem gormModels.XiaomiPlant
			json.Unmarshal(content, &plantItem)

			plantItem.DisplayPid, err = translator.Translate(ctx, plantItem.DisplayPid, apiKey)
			if err != nil {
				return err
			}
			plantItem.Basic.Origin, err = translator.Translate(ctx, plantItem.Basic.Origin, apiKey)
			if err != nil {
				return err
			}
			plantItem.Basic.Production, err = translator.Translate(
				ctx,
				plantItem.Basic.Production,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Basic.Category, err = translator.Translate(
				ctx,
				plantItem.Basic.Category,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Basic.Blooming, err = translator.Translate(
				ctx,
				plantItem.Basic.Blooming,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Basic.Color, err = translator.Translate(ctx, plantItem.Basic.Color, apiKey)
			if err != nil {
				return err
			}

			plantItem.Maintenance.Size, err = translator.Translate(
				ctx,
				plantItem.Maintenance.Size,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Maintenance.Soil, err = translator.Translate(
				ctx,
				plantItem.Maintenance.Soil,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Maintenance.Sunlight, err = translator.Translate(
				ctx,
				plantItem.Maintenance.Sunlight,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Maintenance.Watering, err = translator.Translate(
				ctx,
				plantItem.Maintenance.Watering,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Maintenance.Fertilization, err = translator.Translate(
				ctx,
				plantItem.Maintenance.Fertilization,
				apiKey,
			)
			if err != nil {
				return err
			}
			plantItem.Maintenance.Pruning, err = translator.Translate(
				ctx,
				plantItem.Maintenance.Pruning,
				apiKey,
			)
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

			log.Printf("written %d file of %d", index+1, filesNum)

			return nil
		})
		time.Sleep(300 * time.Millisecond)
	}

	return eg.Wait()
}
