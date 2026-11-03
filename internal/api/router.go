package api

import (
	"net/http"

	"github.com/aaksant/warehouse-sim-be/internal/domains"
	"github.com/aaksant/warehouse-sim-be/internal/simulation"
	"github.com/gin-gonic/gin"
)

func NewRouter(svc *simulation.Service) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/warehouses", listWarehousesHandler(svc))

	api := r.Group("/api")
	api.POST("/simulations/check", checkHandler(svc))

	return r
}

func checkHandler(svc *simulation.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CheckRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := svc.CheckAll(c.Request.Context(), domains.Load{
			TotalTon:     req.TotalTon,
			TonPerPallet: req.TonPerPallet,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Convert domain.FitResult to API DTO
		apiRes := make([]WarehouseCheckResult, len(res))
		for i, r := range res {
			apiRes[i] = WarehouseCheckResult{
				WarehouseID:   r.WarehouseID,
				WarehouseName: r.WarehouseName,
				Fits:          r.Result.Fits,
				PalletsNeeded: r.Result.PalletsNeeded,
				VolumeNeeded:  r.Result.VolumeNeededM3,
				UsedVolume:    r.Result.UsedVolumeM3,
				FreeVolume:    r.Result.FreeVolumeM3,
				CapacityM3:    r.Result.CapacityM3,
			}
		}
		c.JSON(http.StatusOK, gin.H{"results": apiRes})
	}
}

func listWarehousesHandler(svc *simulation.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		whs, err := svc.ListWarehousesWithStatus(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		warehouses := make([]WarehouseCheckResult, len(whs))
		for i, r := range whs {
			warehouses[i] = WarehouseCheckResult{
				WarehouseID:   r.WarehouseID,
				WarehouseName: r.WarehouseName,
				Fits:          r.Fits,
				PalletsNeeded: r.PalletsNeeded,
				VolumeNeeded:  r.VolumeNeededM3,
				UsedVolume:    r.UsedVolumeM3,
				FreeVolume:    r.FreeVolumeM3,
				CapacityM3:    r.CapacityM3,
			}
		}
		c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
	}
}
