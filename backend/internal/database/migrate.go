// Copyright 2026. Kimjibeom. All rights reserved.
package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations automatically reads and executes all .sql files in the migrations directory.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		// Return silently if directory doesn't exist, it might be running in a different context
		log.Printf("Migrations directory '%s' not found, skipping auto-migration", migrationsDir)
		return nil
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsDir, f.Name()))
		}
	}

	// Ensure files are executed in alphabetical order
	sort.Strings(sqlFiles)

	for _, file := range sqlFiles {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		// Execute the migration script
		// pgxpool Exec handles multiple statements natively
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			log.Printf("Error executing migration %s: %v", file, err)
			return err
		}
		log.Printf("Successfully executed migration: %s", filepath.Base(file))
	}

	return nil
}
