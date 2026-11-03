package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaksant/warehouse-sim-be/internal/domains"
	"github.com/aaksant/warehouse-sim-be/internal/simulation"
	"github.com/aaksant/warehouse-sim-be/internal/store/memory"
)

func TestCheckHandler(t *testing.T) {
	store := memory.New()
	service := simulation.NewService(store)

	tests := []struct {
		name       string
		reqBody    CheckRequest
		wantStatus int
		wantFits   bool
	}{
		{
			name:       "load kecil -> 200 fits true di semua gudang",
			reqBody:    CheckRequest{TotalTon: 1, TonPerPallet: 1},
			wantStatus: http.StatusOK,
			wantFits:   true,
		},
		{
			name:       "load sangat besar -> fits false",
			reqBody:    CheckRequest{TotalTon: 10000, TonPerPallet: 1},
			wantStatus: http.StatusOK,
			wantFits:   false,
		},
		{
			name:       "total_ton 0 -> 400",
			reqBody:    CheckRequest{TotalTon: 0, TonPerPallet: 1},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "boundary 268 ton @ 1 ton/pallet fits",
			reqBody:    CheckRequest{TotalTon: 268, TonPerPallet: 1},
			wantStatus: http.StatusOK,
			wantFits:   true,
		},
		{
			name:       "boundary 269 ton @ 1 ton/pallet tidak fits",
			reqBody:    CheckRequest{TotalTon: 269, TonPerPallet: 1},
			wantStatus: http.StatusOK,
			wantFits:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/simulations/check", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := NewRouter(service)
			router.ServeHTTP(w, req)

			if w.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, test.wantStatus)
			}

			if test.wantStatus == http.StatusOK {
				var res map[string]any
				json.Unmarshal(w.Body.Bytes(), &res)

				results := res["results"].([]any)
				if len(results) != 4 {
					t.Fatalf("expected 4 warehouses, got %d", len(results))
				}

				for i, r := range results {
					m := r.(map[string]any)
					if m["fits"] != test.wantFits {
						t.Errorf("warehouse %d fits = %v, want %v", i+1, m["fits"], test.wantFits)
					}
				}
			}
		})
	}
}

func TestCheckHandler_BoundaryCalculation(t *testing.T) {
	spec := domains.DefaultPalletSpec()
	vol := spec.LengthM * spec.WidthM * (spec.LoadHeightPerTonM * 1.0)
	t.Logf("volume per pallet = %.6f m3", vol)

	// 500 m3 / 1.859 m3 per pallet = 269.009...
	// 268 pallets = 498.212 m3 fits
	// 269 pallets = 500.071 m3 does NOT fit
	// 270 pallets = 501.930 m3 does not fit

	pallets268 := 268
	vol268 := float64(pallets268) * vol
	t.Logf("268 pallets volume = %.6f m3 (capacity 500)", vol268)

	pallets269 := 269
	vol269 := float64(pallets269) * vol
	t.Logf("269 pallets volume = %.6f m3 (capacity 500)", vol269)

	pallets270 := 270
	vol270 := float64(pallets270) * vol
	t.Logf("270 pallets volume = %.6f m3 (capacity 500)", vol270)

	if vol268 > 500 {
		t.Fatalf("268 pallets (%.3f) should fit in 500 m3 warehouse", vol268)
	}
	if vol269 <= 500 {
		t.Fatalf("269 pallets (%.3f) should NOT fit in 500 m3 warehouse", vol269)
	}
	if vol270 <= 500 {
		t.Fatalf("270 pallets (%.3f) should NOT fit in 500 m3 warehouse", vol270)
	}
}

func TestListWarehousesHandler(t *testing.T) {
	store := memory.New()
	service := simulation.NewService(store)

	req := httptest.NewRequest(http.MethodGet, "/api/warehouses", nil)
	w := httptest.NewRecorder()

	router := NewRouter(service)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var res map[string]any
	json.Unmarshal(w.Body.Bytes(), &res)

	warehouses := res["warehouses"].([]any)
	if len(warehouses) != 4 {
		t.Fatalf("expected 4 warehouses, got %d", len(warehouses))
	}

	// Empty load, all warehouses should have used_volume_m3 = 0

	for i, w := range warehouses {
		m := w.(map[string]any)
		if m["used_volume_m3"] != 0.0 {
			t.Errorf("warehouse %d used_volume_m3 = %v, want 0", i+1, m["used_volume_m3"])
		}
		if m["free_volume_m3"] != 500.0 {
			t.Errorf("warehouse %d free_volume_m3 = %v, want 500", i+1, m["free_volume_m3"])
		}
		if m["capacity_m3"] != 500.0 {
			t.Errorf("warehouse %d capacity_m3 = %v, want 500", i+1, m["capacity_m3"])
		}
		// used + free should equal capacity
		used := m["used_volume_m3"].(float64)
		free := m["free_volume_m3"].(float64)
		cap := m["capacity_m3"].(float64)
		if used+free != cap {
			t.Errorf("warehouse %d: used_volume + free_volume (%.2f) != capacity (%.2f)", i+1, used+free, cap)
		}
	}
}
