package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGConfig struct {
	DSN             string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func DefaultPGConfig() PGConfig {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:1@172.21.208.1:5432/gotax?sslmode=disable"
	}
	return PGConfig{
		DSN:             dsn,
		MaxConns:        10,
		MinConns:        2,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
	}
}

func NewPool(ctx context.Context, cfg PGConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}
	poolCfg.MaxConns = int32(cfg.MaxConns)
	poolCfg.MinConns = int32(cfg.MinConns)
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	migrations := []string{
		filepath.Join(baseDir, "..", "..", "migrations", "002_gl_schema_circular99.sql"),
		filepath.Join(baseDir, "..", "..", "migrations", "003_company_schema.sql"),
		filepath.Join(baseDir, "..", "..", "migrations", "003_cash_schema.sql"),
		filepath.Join(baseDir, "..", "..", "migrations", "004_advance_schema.sql"),
	filepath.Join(baseDir, "..", "..", "migrations", "004_bank_module.sql"),
	filepath.Join(baseDir, "..", "..", "migrations", "005_purchase_schema.sql"),
	filepath.Join(baseDir, "..", "..", "migrations", "006_sale_schema.sql"),
	filepath.Join(baseDir, "..", "..", "migrations", "007_warehouse_schema.sql"),
}

	for _, path := range migrations {
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", path, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("exec migration %s: %w", path, err)
		}
		log.Printf("migration applied: %s", path)
	}
	return nil
}
