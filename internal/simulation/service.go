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
