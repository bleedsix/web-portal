package cli

import (
	"database/sql"
	"fmt"

	"github.com/bleedsix/web-portal/internal/assets"
	"github.com/bleedsix/web-portal/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // Реєструємо драйвер pgx для database/sql
	migrate "github.com/rubenv/sql-migrate"
)

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: assets.Migrations,
	Root:       "migrations",
}

// openDBForMigrations відкриває класичне *sql.DB з'єднання, яке потрібне для бібліотеки міграцій
func openDBForMigrations(url string) (*sql.DB, error) {
	return sql.Open("pgx", url)
}

func MigrateUp(cfg *config.Config) error {
	// Спеціально відкриваємо з'єднання для мігратора
	db, err := openDBForMigrations(cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to open db for migrations: %w", err)
	}
	defer db.Close()

	count, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	cfg.Log().Info("Migrations applied (UP)", "count", count)
	return nil
}

func MigrateDown(cfg *config.Config) error {
	db, err := openDBForMigrations(cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to open db for migrations: %w", err)
	}
	defer db.Close()

	count, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	cfg.Log().Info("Migrations applied (DOWN)", "count", count)
	return nil
}