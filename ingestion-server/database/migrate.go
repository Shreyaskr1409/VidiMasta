package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(db *pgxpool.Pool, l *log.Logger, ctx context.Context) error {
	l.Println("Running migrations...")
	migrationsPath := "database/migrations"
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory; %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migrations %s: %w", file.Name(), err)
		}

		upPart := string(content)
		if _, err := db.Exec(ctx, upPart); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}

		l.Printf("Applied migration: %s", file.Name())
	}
	l.Println("Migrations ran successfully...")

	return nil
}
