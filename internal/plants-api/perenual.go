package plants_api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/carlmjohnson/requests"
	"github.com/pkg/errors"

	"github.com/cyber_bed/internal/api/convert"
	"github.com/cyber_bed/internal/domain"
	"github.com/cyber_bed/internal/models"
)

type PerenualAPI struct {
	baseURL *url.URL
}

func NewPerenualAPI(url *url.URL, token string) domain.PlantsAPI {
	url.Query().Set("key", token)
	q := url.Query()
	q.Set("key", token)

	url.RawQuery = q.Encode()
	return &PerenualAPI{
		baseURL: url,
	}
}

func clearPremiumItems(plantsResponse models.PerenualsPlantResponse) models.PerenualsPlantResponse {
	var noPremiumPlants models.PerenualsPlantResponse
	noPremiumPlants.Data = make([]models.PerenualPlant, 0)

	for _, plant := range plantsResponse.Data {
		if plant.Cycle != "Upgrade Plans To Premium/Supreme - https://www.perenual.com/subscription-api-pricing. I'm sorry" {
			noPremiumPlants.Data = append(noPremiumPlants.Data, plant)
		}
	}

	return noPremiumPlants
}

func (p *PerenualAPI) SearchByName(ctx context.Context, name string) ([]models.Plant, error) {
	u := p.baseURL
	q := u.Query()
	q.Set("q", name)

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("species-list")

	var resp models.PerenualsPlantResponse

	if err := requests.
		URL(apiURL.String()).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to search plant by name")
	}

	resp = clearPremiumItems(resp)
	return convert.InputSearchPerenaulResultsToModels(resp, 30), nil
}

func (p *PerenualAPI) SearchByID(ctx context.Context, id uint64) (models.Plant, error) {
	u := p.baseURL
	q := u.Query()

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("species").JoinPath("details").JoinPath(strconv.FormatUint(id, 10))

	var resp models.PerenualPlant

	if err := requests.
		URL(apiURL.String()).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return models.Plant{}, errors.Wrap(err, "failed to search plant by name")
	}
	return convert.SearchItemToPlantModel(resp), nil
}

func (p *PerenualAPI) GetPage(ctx context.Context, pageNum uint64) ([]models.Plant, error) {
	u := p.baseURL
	q := u.Query()
	q.Set("page", strconv.FormatUint(pageNum, 10))

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("species-list")

	var resp models.PerenualsPlantResponse

	if err := requests.
		URL(apiURL.String()).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to search plant by name")
	}

	resp = clearPremiumItems(resp)
	return convert.InputSearchPerenaulResultsToModels(resp, 30), nil
}
