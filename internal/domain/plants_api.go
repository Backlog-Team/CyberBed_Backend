package domain

import (
	"context"

	"github.com/cyber_bed/internal/models"
)

type PlantsAPI interface {
	SearchByName(ctx context.Context, name string) ([]models.Plant, error)
	SearchByID(ctx context.Context, id uint64) (models.Plant, error)
	GetPage(ctx context.Context, pageNum uint64) ([]models.Plant, error)
}
