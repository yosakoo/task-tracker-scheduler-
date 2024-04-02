package main

import (
	"log"

	"github.com/yosakoo/task-tracker-scheduler-/internal/config"
	"github.com/yosakoo/task-tracker-scheduler-/internal/app"
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