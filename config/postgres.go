package config

import (
	"context"
	"log"
	"os"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitPostgres() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return fmt.Errorf("DATABASE_URL is missing in .env")
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return err
	}

	DB = conn
	log.Println("✅ PostgreSQL connected")
	return nil
}
