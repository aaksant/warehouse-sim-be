package simulation

import (
	"context"

	"github.com/aaksant/warehouse-sim-be/internal/domains"
	"github.com/aaksant/warehouse-sim-be/internal/store"
)

type Service struct {
	store store.Store
	spec  domains.PalletSpec
}

func NewService(s store.Store) *Service {
	return &Service{
		store: s,
		spec:  domains.DefaultPalletSpec(),
	}
}

// GetWarehouse returns a warehouse by ID
func (svc *Service) GetWarehouse(ctx context.Context, warehouseID int) (domains.Warehouse, error) {
	warehouses, err := svc.store.ListWarehouses(ctx)
	if err != nil {
		return domains.Warehouse{}, err
	}
	for _, w := range warehouses {
		if w.ID == warehouseID {
			return w, nil
		}
	}
	return domains.Warehouse{}, domains.ErrNotImplemented
}

// GetUsedVolume returns the used volume for a warehouse based on its batches
func (svc *Service) GetUsedVolume(ctx context.Context, warehouseID int) (float64, error) {
	batches, err := svc.store.ListBatches(ctx, warehouseID)
	if err != nil {
		return 0, err
	}
	usedM3 := 0.0
	for _, b := range batches {
		pallets, err := domains.PalletsNeeded(domains.Load{
			TotalTon:     b.TotalTon,
			TonPerPallet: b.TonPerPallet,
		})
		if err != nil {
			return 0, err
		}
		if pallets > 0 {
			height := svc.spec.LoadHeightPerTonM * b.TonPerPallet
			vol := svc.spec.LengthM * svc.spec.WidthM * height
			usedM3 += float64(pallets) * vol
		}
	}
	return usedM3, nil
}

// SaveBatch saves a batch to the store
func (svc *Service) SaveBatch(ctx context.Context, b domains.Batch) error {
	return svc.store.SaveBatch(ctx, b)
}

// FindBatch finds a batch by surat_jalan, material, and received_at
func (svc *Service) FindBatch(ctx context.Context, suratJalan, material, receivedAt string) (*domains.Batch, error) {
	return svc.store.FindBatch(ctx, suratJalan, material, receivedAt)
}

// GetSpec returns the pallet spec
func (svc *Service) GetSpec() domains.PalletSpec {
	return svc.spec
}

// WarehouseResult adalah hasil simulasi untuk satu gudang.
type WarehouseResult struct {
	WarehouseID   int               `json:"warehouse_id"`
	WarehouseName string            `json:"warehouse_name"`
	Result        domains.FitResult `json:"result"`
}

// WarehouseStatus adalah status satu gudang (used/free volume).
type WarehouseStatus struct {
	WarehouseID    int     `json:"warehouse_id"`
	WarehouseName  string  `json:"warehouse_name"`
	Fits           bool    `json:"fits"`
	PalletsNeeded  int     `json:"pallets_needed"`
	VolumeNeededM3 float64 `json:"volume_needed_m3"`
	UsedVolumeM3   float64 `json:"used_volume_m3"`
	FreeVolumeM3   float64 `json:"free_volume_m3"`
	CapacityM3     float64 `json:"capacity_m3"`
}

// ListWarehousesWithStatus mengembalikan status setiap gudang (used/free volume).
func (svc *Service) ListWarehousesWithStatus(
	ctx context.Context,
) ([]WarehouseStatus, error) {
	warehouses, err := svc.store.ListWarehouses(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]WarehouseStatus, 0, len(warehouses))
	for _, w := range warehouses {
		batches, err := svc.store.ListBatches(ctx, w.ID)
		if err != nil {
			return nil, err
		}
		usedM3 := 0.0
		for _, b := range batches {
			pallets, err := domains.PalletsNeeded(domains.Load{
				TotalTon:     b.TotalTon,
				TonPerPallet: b.TonPerPallet,
			})
			if err != nil {
				return nil, err
			}
			if pallets > 0 {
				height := svc.spec.LoadHeightPerTonM * b.TonPerPallet
				vol := svc.spec.LengthM * svc.spec.WidthM * height
				usedM3 += float64(pallets) * vol
			}
		}
		out = append(out, WarehouseStatus{
			WarehouseID:    w.ID,
			WarehouseName:  w.Name,
			Fits:           true, // N/A for warehouse status
			PalletsNeeded:  0,    // N/A
			VolumeNeededM3: 0,    // N/A
			UsedVolumeM3:   usedM3,
			FreeVolumeM3:   w.CapacityM3 - usedM3,
			CapacityM3:     w.CapacityM3,
		})
	}
	return out, nil
}

// CheckAll menyimulasikan load terhadap keempat gudang sekaligus.
func (svc *Service) CheckAll(
	ctx context.Context,
	l domains.Load,
) ([]WarehouseResult, error) {
	warehouses, err := svc.store.ListWarehouses(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]WarehouseResult, 0, len(warehouses))
	for _, w := range warehouses {
		batches, err := svc.store.ListBatches(ctx, w.ID)
		if err != nil {
			return nil, err
		}
		usedM3 := 0.0
		for _, b := range batches {
			pallets, err := domains.PalletsNeeded(domains.Load{
				TotalTon:     b.TotalTon,
				TonPerPallet: b.TonPerPallet,
			})
			if err != nil {
				return nil, err
			}
			if pallets > 0 {
				// Reuse volumePerPalletM3 logic
				height := svc.spec.LoadHeightPerTonM * b.TonPerPallet
				vol := svc.spec.LengthM * svc.spec.WidthM * height
				usedM3 += float64(pallets) * vol
			}
		}
		res, err := domains.CanFit(w, usedM3, l, svc.spec)
		if err != nil {
			return nil, err
		}
		out = append(out, WarehouseResult{
			WarehouseID:   w.ID,
			WarehouseName: w.Name,
			Result:        res,
		})

	}
	return out, nil
}
