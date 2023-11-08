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

type TrefleAPI struct {
	baseURL      *url.URL
	countResults int
}

func NewTrefleAPI(baseURL *url.URL, countResults int, token string) domain.PlantsAPI {
	baseURL.Query().Set("token", token)
	q := baseURL.Query()
	q.Set("token", token)

	baseURL.RawQuery = q.Encode()

	return &TrefleAPI{
		baseURL:      baseURL,
		countResults: countResults,
	}
}

func (t *TrefleAPI) SearchByName(
	ctx context.Context,
	name string,
) ([]models.Plant, error) {
	u := t.baseURL
	q := u.Query()
	q.Set("q", name)

	u.RawQuery = q.Encode()
	apiURL := u.JoinPath("search")

	var resp models.SearchSliceResponse

	if err := requests.
		URL(apiURL.String()).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to search plant by name")
	}

	return convert.InputSearchTrefleResultsToModels(resp, t.countResults), nil
}

func (t *TrefleAPI) SearchByID(ctx context.Context, id uint64) (models.Plant, error) {
	u := t.baseURL
	q := u.Query()
	u.RawQuery = q.Encode()
	apiURL := u.JoinPath(strconv.FormatUint(id, 10)).String()

	var resp models.SearchResponse
	if err := requests.
		URL(apiURL).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return models.Plant{}, errors.Wrap(err, "failed to search plant by id")
	}
	return convert.SearchTrefleItemToPlantModel(resp.Data), nil
}

func (t *TrefleAPI) GetPage(ctx context.Context, pageNum uint64) ([]models.Plant, error) {
	u := t.baseURL
	q := u.Query()
	q.Set("page", strconv.FormatUint(pageNum, 10))

	u.RawQuery = q.Encode()

	var resp models.SearchSliceResponse
	if err := requests.
		URL(u.String()).
		Method(http.MethodGet).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to get page")
	}
	return convert.InputSearchTrefleResultsToModels(resp, t.countResults), nil
}
