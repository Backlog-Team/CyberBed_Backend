package usecase

import (
	"context"
	"mime/multipart"

	"github.com/pkg/errors"

	"github.com/cyber_bed/internal/domain"
	httpModels "github.com/cyber_bed/internal/models/http"
	domainRecognition "github.com/cyber_bed/internal/recognize-api"
)

type usecase struct {
	apiRecognition domainRecognition.API
	apiPlants      domain.PlantsAPI
	plantsUsecase  domain.PlantsUsecase
}

func New(
	api domainRecognition.API,
	apiPlants domain.PlantsAPI,
	plantsUsecase domain.PlantsUsecase,
) domainRecognition.Usecase {
	return usecase{
		apiRecognition: api,
		apiPlants:      apiPlants,
		plantsUsecase:  plantsUsecase,
	}
}

func (u usecase) Recognize(
	ctx context.Context,
	formdata *multipart.Form,
	project string,
) ([]httpModels.XiaomiPlant, error) {
	recognized, err := u.apiRecognition.Recognize(ctx, formdata, httpModels.Project(project))
	if err != nil {
		return nil, errors.Wrap(err, "failed to recognize images")
	}

	plants := make([]httpModels.XiaomiPlant, 0)
	curIndex := 0
	for _, plant := range recognized {
		found, err := u.plantsUsecase.GetPlantByName(plant.CommonName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to search plant")
		}

		plants = append(plants, found...)
		for range found {
			plants[curIndex].PredictionScore = plant.PredictionScore
			curIndex += 1
		}
	}

	return plants, nil
}
