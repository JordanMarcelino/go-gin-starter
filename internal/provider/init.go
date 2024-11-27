package provider

import (
	"database/sql"

	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/internal/database"
)

var (
	db *sql.DB
)

func InitGlobal(cfg *config.Config) {
	db = database.InitPostgres(cfg)
}
