package bids

import (
	"context"

	"go-metricscol/internal/models"
)

type UseCase interface {
	Find(ctx context.Context, name string, mType models.MetricType) (*models.Metric, error)
	Update(ctx context.Context, metric models.Metric) error
	Updates(ctx context.Context, metrics []models.Metric) error
	GetAll(ctx context.Context) ([]models.Metric, error)
}
