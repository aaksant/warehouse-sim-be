package main

import (
	"log"

	"github.com/aaksant/warehouse-sim-be/internal/api"
	"github.com/aaksant/warehouse-sim-be/internal/config"
	"github.com/aaksant/warehouse-sim-be/internal/simulation"
	"github.com/aaksant/warehouse-sim-be/internal/store/memory"
)

func main() {
	config := config.Load()
	store := memory.New()
	service := simulation.NewService(store)
	router := api.NewRouter(service)

	log.Printf("listening on %s", config.Addr)
	if err := router.Run(config.Addr); err != nil {
		log.Fatal(err)
	}
}
