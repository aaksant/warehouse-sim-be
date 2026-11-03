package domains

import "time"

// Batch adalah barang dengan surat jalan + bahan yang sama di hari yang sama.
type Batch struct {
	ID           string
	SuratJalan   string
	Material     string
	ReceivedAt   time.Time
	TotalTon     float64
	TonPerPallet float64
	WarehouseID  int
}
