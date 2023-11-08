package domain

import (
	"context"
	"mime/multipart"

	"github.com/labstack/echo/v4"

	"github.com/cyber_bed/internal/models"
)

type API interface {
	Recognize(
		ctx context.Context,
		formdata *multipart.Form,
		project models.Project,
	) ([]models.Plant, error)
}

type Usecase interface {
	Recognize(
		ctx context.Context,
		formdata *multipart.Form,
		project string,
	) ([]models.XiaomiPlant, error)
}

type Handler interface {
	Recognize(c echo.Context) error
}
