package convert

import httpModels "github.com/cyber_bed/internal/models/http"

func InputRecognitionResultsToModels(results httpModels.RecResponse, count int) []httpModels.Plant {
	plants := make([]httpModels.Plant, 0)
	counter := 0

	for _, result := range results.Results {
		if counter+1 == count {
			break
		}

		plants = append(plants, httpModels.Plant{
			CommonName: result.Species.Name,
		})
		counter++
	}

	return plants
}

func InputSearchTrefleResultsToModels(
	results httpModels.SearchSliceResponse,
	count int,
) []httpModels.Plant {
	plants := make([]httpModels.Plant, 0)
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
	results httpModels.PerenualsPlantResponse,
	count int,
) []httpModels.Plant {
	plants := make([]httpModels.Plant, 0)
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

func SearchTrefleItemToPlantModel(res httpModels.ItemPlantResponse) httpModels.Plant {
	return httpModels.Plant{
		ID:             uint64(res.ID),
		ScientificName: []string{res.ScName},
		ImageUrl:       res.ImageURL,
	}
}

func SearchItemToPlantModel(res httpModels.PerenualPlant) httpModels.Plant {
	return httpModels.Plant{
		ID:             uint64(res.ID),
		CommonName:     res.CommonName,
		ImageUrl:       res.ImageURL.URL,
		ScientificName: res.ScientificName,
		OtherName:      res.OtherName,
		Cycle:          res.Cycle,
		Watering:       res.Watering,
	}
}
