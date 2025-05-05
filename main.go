package main

import (
	"effectiveMobile/api"
	"effectiveMobile/env"
	"effectiveMobile/internal/db"
	"effectiveMobile/internal/utils"
)

func main() {
	cfg, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	utils.Initialize(cfg.DebugMode)

	rep, err := db.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	err = db.ApplyMigrations(cfg)
	if err != nil {
		panic(err)
	}

	external := api.NewExternal()
	handlers := api.NewHandlers(rep, external)
	err = handlers.Start(cfg)
	if err != nil {
		panic(err)
	}
}
