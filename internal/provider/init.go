package provider

import (
	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/internal/database"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func InitGlobal(cfg *config.Config) {
	db = database.InitPostgres(cfg)
}
