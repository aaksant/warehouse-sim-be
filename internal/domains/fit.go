package domains

import (
	"errors"
	"math"
)

var ErrNotImplemented = errors.New("not implemented")

// Load adalah muatan yang ingin disimulasikan.
type Load struct {
	TotalTon     float64
	TonPerPallet float64 // 1.0, 1.375, dst. (ada varian berat)
}

// FitResult adalah hasil simulasi untuk satu gudang.
type FitResult struct {
	Fits           bool    `json:"fits"`
	PalletsNeeded  int     `json:"pallets_needed"`
	VolumeNeededM3 float64 `json:"volume_needed_m3"`
	UsedVolumeM3   float64 `json:"used_volume_m3"`
	FreeVolumeM3   float64 `json:"free_volume_m3"`
	CapacityM3     float64 `json:"capacity_m3"`
}

// PalletsNeeded menghitung jumlah pallet untuk sebuah load (dibulatkan ke atas).
func PalletsNeeded(l Load) (int, error) {
	if l.TotalTon <= 0 || l.TonPerPallet <= 0 {
		return 0, errors.New("total_ton and ton_per_pallet must be positive")
	}
	// Kurangi epsilon kecil agar hasil seperti 3.0000000001 tidak naik jadi 4.
	return int(math.Ceil(l.TotalTon/l.TonPerPallet - 1e-9)), nil
}

// CanFit menentukan apakah load masih muat di gudang w, diberi volume yang
// sudah terpakai (usedM3). "Muat" dihitung berdasarkan volume.
func CanFit(
	w Warehouse,
	usedM3 float64,
	l Load, spec PalletSpec,
) (FitResult, error) {
	pallets, err := PalletsNeeded(l)
	if err != nil {
		return FitResult{}, err
	}

	need := float64(pallets) * volumePerPalletM3(l, spec)
	free := math.Max(w.CapacityM3-usedM3, 0)

	return FitResult{
		Fits:           need <= free,
		PalletsNeeded:  pallets,
		VolumeNeededM3: need,
		UsedVolumeM3:   usedM3,
		FreeVolumeM3:   free,
		CapacityM3:     w.CapacityM3,
	}, nil
}

// TODO: asumsi A, tinggi naik sebanding dengan berat per pallet.
// volumePerPalletM3 menghitung volume satu pallet berisi muatan.
// Tinggi 130 cm per 1 ton sudah termasuk pallet, jadi tidak ditambah lagi.
func volumePerPalletM3(l Load, spec PalletSpec) float64 {
	height := spec.LoadHeightPerTonM * l.TonPerPallet
	return spec.LengthM * spec.WidthM * height
}
