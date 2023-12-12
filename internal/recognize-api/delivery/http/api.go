package http

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	httpModels "github.com/cyber_bed/internal/models/http"
	domain "github.com/cyber_bed/internal/recognize-api"
)

type RecognitionHandler struct {
	usecase domain.Usecase
}

func NewHandler(usecase domain.Usecase) domain.Handler {
	return &RecognitionHandler{
		usecase: usecase,
	}
}

func (r *RecognitionHandler) Recognize(c echo.Context) error {
	formdata, err := c.MultipartForm()
	if err != nil {
		return errors.Wrap(err, "failed to export formdata")
	}

	recognize, err := r.usecase.Recognize(
		c.Request().Context(),
		formdata,
		string(httpModels.AllProject),
	)
	if err != nil {
		return errors.Wrap(err, "failed to recognize plant")
	}

	sort.SliceStable(recognize, func(i, j int) bool {
		return recognize[i].PredictionScore >= recognize[j].PredictionScore
	})

	return c.JSON(http.StatusOK, recognize)
}
