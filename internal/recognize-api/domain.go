package domain

import (
	"context"
	"mime/multipart"

	"github.com/labstack/echo/v4"

	httpModels "github.com/cyber_bed/internal/models/http"
)

type API interface {
	Recognize(
		ctx context.Context,
		formdata *multipart.Form,
		project httpModels.Project,
	) ([]httpModels.Plant, error)
}

type Usecase interface {
	Recognize(
		ctx context.Context,
		formdata *multipart.Form,
		project string,
	) ([]httpModels.XiaomiPlant, error)
}

type Handler interface {
	Recognize(c echo.Context) error
}
