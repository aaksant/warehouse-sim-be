package memory

import (
	"context"
	"sync"

	"github.com/aaksant/warehouse-sim-be/internal/domains"
)

type Store struct {
	mu         sync.RWMutex
	warehouses []domains.Warehouse
	batches    []domains.Batch
}

func New() *Store {
	return &Store{
		warehouses: []domains.Warehouse{
			{ID: 1, Name: "Gudang 1", CapacityM3: 500}, // isi kapasitas kalo udah dapet layout
			{ID: 2, Name: "Gudang 2", CapacityM3: 500},
			{ID: 3, Name: "Gudang 3", CapacityM3: 500},
			{ID: 4, Name: "Gudang 4", CapacityM3: 500},
		},
	}
}

func (s *Store) ListWarehouses(ctx context.Context) ([]domains.Warehouse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domains.Warehouse, len(s.warehouses))
	copy(out, s.warehouses)
	return out, nil
}

func (s *Store) ListBatches(
	ctx context.Context,
	warehouseID int,
) ([]domains.Batch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []domains.Batch
	for _, b := range s.batches {
		if b.WarehouseID == warehouseID {
			out = append(out, b)
		}
	}
	return out, nil
}

func (s *Store) SaveBatch(ctx context.Context, b domains.Batch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches = append(s.batches, b)
	return nil
}
