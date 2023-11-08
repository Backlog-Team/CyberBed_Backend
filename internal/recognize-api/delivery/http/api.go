package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/cyber_bed/internal/models"
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
		string(models.AllProject),
	)
	if err != nil {
		return errors.Wrap(err, "failed to recognize plant")
	}

	return c.JSON(http.StatusOK, recognize)
}
