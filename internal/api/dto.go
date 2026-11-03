package api

type CheckRequest struct {
	TotalTon     float64 `json:"total_ton" binding:"required,gt=0"`
	TonPerPallet float64 `json:"ton_per_pallet" binding:"required,gt=0"`
}

type CheckResponse struct {
	Results []WarehouseCheckResult `json:"results"`
}

type WarehouseCheckResult struct {
	WarehouseID   int     `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	Fits          bool    `json:"fits"`
	PalletsNeeded int     `json:"pallets_needed"`
	VolumeNeeded  float64 `json:"volume_needed_m3"`
	UsedVolume    float64 `json:"used_volume_m3"`
	FreeVolume    float64 `json:"free_volume_m3"`
	CapacityM3    float64 `json:"capacity_m3"`
}
