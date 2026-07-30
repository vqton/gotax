package db

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunGolangMigrate(dsn string) error {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	migrationsPath := filepath.Join(baseDir, "..", "..", "migrations")

	sourceURL := "file://" + migrationsPath
	databaseURL := dsn
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		databaseURL = strings.TrimPrefix(databaseURL, prefix)
	}
	databaseURL = "pgx5://" + databaseURL

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	zap.L().Info("golang-migrate: migrations complete")
	return nil
}
