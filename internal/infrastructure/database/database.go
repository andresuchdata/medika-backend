package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"medika-backend/internal/infrastructure/config"
	"medika-backend/internal/infrastructure/persistence/models"
)

func New(cfg config.DatabaseConfig) (*bun.DB, error) {
	// Create PostgreSQL connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.URL)))

	// Configure connection pool
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqldb.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Create Bun DB instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query hook for debugging (development only)
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Register models
	registerModels(db)

	return db, nil
}

func registerModels(db *bun.DB) {
	db.RegisterModel(
		(*models.User)(nil),
		(*models.UserProfile)(nil),
		(*models.Organization)(nil),
		(*models.Appointment)(nil),
		(*models.Room)(nil),
		(*models.PatientQueue)(nil),
		(*models.Notification)(nil),
		(*models.Media)(nil),
	)
}

func Migrate(db *bun.DB) error {
	ctx := context.Background()

	// Read and execute migration files
	migrationPath := "migrations"
	files, err := filepath.Glob(filepath.Join(migrationPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		fmt.Printf("✅ Applied migration: %s\n", filepath.Base(file))
	}

	return nil
}
