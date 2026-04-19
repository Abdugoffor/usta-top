package config

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations() {
	ctx := context.Background()

	// Tracking jadvali — bir marta yaratiladi, hech qachon o'chirilmaydi
	_, err := DB.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id         SERIAL       PRIMARY KEY,
			name       VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		)
	`)
	
	if err != nil {
		log.Fatal("❌ schema_migrations table error:", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	{
		if err != nil {
			log.Fatal("❌ migrations dir read error:", err)
		}
	}

	// Fayl nomiga ko'ra tartiblash (001_, 002_, ...)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	applied := 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		name := entry.Name()

		// Allaqachon ishlatiganmi?
		var count int
		{
			if err := DB.QueryRow(ctx,
				`SELECT COUNT(*) FROM schema_migrations WHERE name = $1`, name,
			).Scan(&count); err != nil {
				log.Fatalf("❌ Migration check error [%s]: %v", name, err)
			}
		}
		if count > 0 {
			continue
		}

		// SQL o'qi va bajar
		content, err := migrationFiles.ReadFile("migrations/" + name)
		{
			if err != nil {
				log.Fatalf("❌ Migration read error [%s]: %v", name, err)
			}
		}

		if _, err := DB.Exec(ctx, string(content)); err != nil {
			log.Fatalf("❌ Migration failed [%s]: %v", name, err)
		}

		// Bajarilganligini yoz
		if _, err := DB.Exec(ctx,
			`INSERT INTO schema_migrations (name) VALUES ($1)`, name,
		); err != nil {
			log.Fatalf("❌ Migration record error [%s]: %v", name, err)
		}

		log.Printf("✅ Migration applied: %s", name)
		applied++
	}

	if applied == 0 {
		log.Println("✅ Migrations: nothing new")
	} else {
		log.Printf("✅ Migrations: %d applied", applied)
	}
}
