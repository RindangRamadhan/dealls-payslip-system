package main

import (
	"log"

	"github.com/RindangRamadhan/dealls-payslip-system/config"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
