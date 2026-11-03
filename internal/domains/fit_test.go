package domains

import "testing"

func TestPalletsNeeded(t *testing.T) {
	tests := []struct {
		name string
		load Load
		want int
	}{
		{"3 ton @ 1 ton/pallet", Load{TotalTon: 3, TonPerPallet: 1}, 3},
		{
			"5 ton @ 1.375 ton/pallet dibulatkan ke atas",
			Load{
				TotalTon: 5, TonPerPallet: 1.375,
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PalletsNeeded(tt.load)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCanFit_Basic(t *testing.T) {
	spec := DefaultPalletSpec()
	load := Load{TotalTon: 3, TonPerPallet: 1}

	t.Run("gudang kosong dan besar: muat", func(t *testing.T) {
		res, err := CanFit(Warehouse{ID: 1, CapacityM3: 10000}, 0, load, spec)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res.Fits {
			t.Fatal("expected load to fit")
		}
	})

	t.Run("kapasitas nol: tidak muat", func(t *testing.T) {
		res, err := CanFit(Warehouse{ID: 1, CapacityM3: 0}, 0, load, spec)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Fits {
			t.Fatal("expected load not to fit")
		}
	})
}

func TestCanFit_Boundary(t *testing.T) {
	spec := DefaultPalletSpec()
	load := Load{TotalTon: 1, TonPerPallet: 1} // 1 pallet = 1,859 m³

	res, _ := CanFit(Warehouse{CapacityM3: 1.85}, 0, load, spec)
	if res.Fits {
		t.Fatal("1.85 m3 should not fit 1.859 m3")
	}
	res, _ = CanFit(Warehouse{CapacityM3: 1.87}, 0, load, spec)
	if !res.Fits {
		t.Fatal("1.87 m3 should fit 1.859 m3")
	}
}
