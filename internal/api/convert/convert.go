package convert

import "github.com/cyber_bed/internal/models"

func InputRecognitionResultsToModels(results models.RecResponse, count int) []models.Plant {
	plants := make([]models.Plant, 0)
	counter := 0

	for _, result := range results.Results {
		if counter+1 == count {
			break
		}

		plants = append(plants, models.Plant{
			CommonName: result.Species.Name,
		})
		counter++
	}

	return plants
}

func InputSearchTrefleResultsToModels(
	results models.SearchSliceResponse,
	count int,
) []models.Plant {
	plants := make([]models.Plant, 0)
	counter := 0

	for _, result := range results.Data {
		if counter+1 == count {
			break
		}

		plants = append(plants, SearchTrefleItemToPlantModel(result))
		counter++
	}

	return plants
}

func InputSearchPerenaulResultsToModels(
	results models.PerenualsPlantResponse,
	count int,
) []models.Plant {
	plants := make([]models.Plant, 0)
	counter := 0

	for _, result := range results.Data {
		if counter+1 == count {
			break
		}

		plants = append(plants, SearchItemToPlantModel(result))
		counter++
	}

	return plants
}

func SearchTrefleItemToPlantModel(res models.ItemPlantResponse) models.Plant {
	return models.Plant{
		ID:             uint64(res.ID),
		ScientificName: []string{res.ScName},
		ImageUrl:       res.ImageURL,
	}
}

func SearchItemToPlantModel(res models.PerenualPlant) models.Plant {
	return models.Plant{
		ID:             uint64(res.ID),
		CommonName:     res.CommonName,
		ImageUrl:       res.ImageURL.URL,
		ScientificName: res.ScientificName,
		OtherName:      res.OtherName,
		Cycle:          res.Cycle,
		Watering:       res.Watering,
		Sunlight:       []interface{}{res.Sunlight},
	}
}
