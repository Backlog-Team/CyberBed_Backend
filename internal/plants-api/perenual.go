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
	httpModels "github.com/cyber_bed/internal/models/http"
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

func clearPremiumItems(
	plantsResponse httpModels.PerenualsPlantResponse,
) httpModels.PerenualsPlantResponse {
	var noPremiumPlants httpModels.PerenualsPlantResponse
	noPremiumPlants.Data = make([]httpModels.PerenualPlant, 0)

	for _, plant := range plantsResponse.Data {
		if plant.Cycle != "Upgrade Plans To Premium/Supreme - https://www.perenual.com/subscription-api-pricing. I'm sorry" {
			noPremiumPlants.Data = append(noPremiumPlants.Data, plant)
		}
	}

	return noPremiumPlants
}

func (p *PerenualAPI) SearchByName(ctx context.Context, name string) ([]httpModels.Plant, error) {
	u := p.baseURL
	q := u.Query()
	q.Set("q", name)

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("species-list")

	var resp httpModels.PerenualsPlantResponse

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

func (p *PerenualAPI) SearchByID(ctx context.Context, id uint64) (httpModels.Plant, error) {
	u := p.baseURL
	q := u.Query()

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("species").JoinPath("details").JoinPath(strconv.FormatUint(id, 10))

	var resp httpModels.PerenualPlant

	if err := requests.
		URL(apiURL.String()).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return httpModels.Plant{}, errors.Wrap(err, "failed to search plant by name")
	}
	return convert.SearchItemToPlantModel(resp), nil
}

func (p *PerenualAPI) GetPage(ctx context.Context, pageNum uint64) ([]httpModels.Plant, error) {
	u := p.baseURL
	q := u.Query()
	q.Set("page", strconv.FormatUint(pageNum, 10))

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("species-list")

	var resp httpModels.PerenualsPlantResponse

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
