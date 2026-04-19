package config

import (
	"auth_service/helper"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func DBConnect() *pgxpool.Pool {
	driver := helper.ENV("DB_DRIVER")

	if driver != "postgres" {
		log.Fatal("❌ pgx faqat PostgreSQL bilan ishlaydi (DB_DRIVER=postgres)")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		helper.ENV("DB_USER"),
		helper.ENV("DB_PASSWORD"),
		helper.ENV("DB_HOST"),
		helper.ENV("DB_PORT"),
		helper.ENV("DB_NAME"),
		helper.ENV("DB_SSLMODE"),
		helper.ENV("DB_TIMEZONE"),
	)

	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("❌ DSN parse error:", err)
	}

	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatal("❌ DB ping error:", err)
	}

	log.Println("✅ Connected to PostgreSQL (pgxpool) 🚀")

	DB = db

	RunMigrations()

	return db
}
