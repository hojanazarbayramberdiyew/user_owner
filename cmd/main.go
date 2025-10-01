package main

import (
	"fmt"
	"user_owner/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

}

func connectDB() {

	cfg := config.GetConfig()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Storage.DbName, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username)

	pgxpool.Pool()
}
