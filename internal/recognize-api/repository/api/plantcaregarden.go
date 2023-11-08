package postgres

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests"
	"github.com/pkg/errors"

	"github.com/cyber_bed/internal/api/convert"
	"github.com/cyber_bed/internal/models"
	domain "github.com/cyber_bed/internal/recognize-api"
)

type RecognitionAPI struct {
	baseURL      *url.URL
	apiKey       string
	imageField   string
	maxImages    int
	countResults int
}

func NewRecognitionAPI(
	url *url.URL,
	apiKey string,
	imageField string,
	maxImages int,
	countResults int,
) domain.API {
	return &RecognitionAPI{
		baseURL:      url,
		apiKey:       apiKey,
		imageField:   imageField,
		maxImages:    maxImages,
		countResults: countResults,
	}
}

func (r *RecognitionAPI) Recognize(
	ctx context.Context,
	formdata *multipart.Form,
	project models.Project,
) ([]models.Plant, error) {
	images, ok := formdata.File[r.imageField]
	if !ok {
		return nil, models.ErrNoImages
	}

	if len(images) > 5 {
		return nil, errors.Wrapf(models.ErrTooManyImages, "required %d", r.maxImages)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	file, _ := images[0].Open()

	fw, _ := w.CreateFormFile(r.imageField, images[0].Filename)

	io.Copy(fw, file)

	w.Close()

	apiURL := r.baseURL.JoinPath(string(project))
	q := apiURL.Query()
	q.Set("api-key", r.apiKey)

	apiURL.RawQuery = q.Encode()

	var resp models.RecResponse

	if err := requests.
		URL(apiURL.String()).
		ContentType(w.FormDataContentType()).
		Method(http.MethodPost).
		BodyReader(&b).
		Header("api-key", r.apiKey).
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to recognize plant by image")
	}

	return convert.InputRecognitionResultsToModels(resp, r.countResults), nil
}
