package domains

// Warehouse adalah satu gudang. Keempat gudang karakteristiknya sama;
// layout internal belum final, jadi kapasitas diisi dari data, bukan hardcode.
type Warehouse struct {
	ID         int
	Name       string
	CapacityM3 float64
}
