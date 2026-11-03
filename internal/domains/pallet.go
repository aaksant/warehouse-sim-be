package domains

// PalletSpec adalah spesifikasi pallet dan aturan tumpuk.
// Nilai default berasal dari info grup dan BELUM dikonfirmasi final.
type PalletSpec struct {
	LengthM float64
	WidthM  float64
	HeightM float64
	// MaxStack adalah jumlah maksimum pallet yang boleh ditumpuk.
	MaxStack int
	// LoadHeightPerTonM adalah tinggi rata-rata per 1 ton muatan (termasuk pallet).
	LoadHeightPerTonM float64
}

func DefaultPalletSpec() PalletSpec {
	return PalletSpec{
		LengthM:           1.30,
		WidthM:            1.10,
		HeightM:           0.06,
		MaxStack:          3,
		LoadHeightPerTonM: 1.30,
	}
}
