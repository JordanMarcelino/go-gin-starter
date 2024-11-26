package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitPostgres(cfg *config.Config) *sql.DB {
	dbCfg := cfg.Database

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		dbCfg.Host,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.DbName,
		dbCfg.Port,
		dbCfg.Sslmode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Fatalf("error initializing database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Log.Fatalf("error connecting to database: %v", err)
	}

	db.SetMaxIdleConns(dbCfg.MaxIdleConn)
	db.SetMaxOpenConns(dbCfg.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(dbCfg.MaxConnLifetime) * time.Minute)

	return db
}
