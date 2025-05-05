package db

import (
	"effectiveMobile/env"
	"effectiveMobile/internal/utils"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(cfg env.Config) error {
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSlMode)

	m, err := migrate.New("file://migrations", dbUrl)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer m.Close()

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		utils.LogInfo("No changes: migration is applied")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error init migrations: %w", err)
	}

	utils.LogInfo("Migrations have been successfully applied")
	return nil
}
