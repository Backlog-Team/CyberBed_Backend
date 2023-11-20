package domain

import (
	"context"

	httpModels "github.com/cyber_bed/internal/models/http"
)

type PlantsAPI interface {
	SearchByName(ctx context.Context, name string) ([]httpModels.Plant, error)
	SearchByID(ctx context.Context, id uint64) (httpModels.Plant, error)
	GetPage(ctx context.Context, pageNum uint64) ([]httpModels.Plant, error)
}
