package store

import (
	"context"

	"github.com/aaksant/warehouse-sim-be/internal/domains"
)

type Store interface {
	ListWarehouses(ctx context.Context) ([]domains.Warehouse, error)
	ListBatches(ctx context.Context, warehouseID int) ([]domains.Batch, error)
	SaveBatch(ctx context.Context, b domains.Batch) error
	FindBatch(ctx context.Context, suratJalan, material, receivedAt string) (*domains.Batch, error)
}
