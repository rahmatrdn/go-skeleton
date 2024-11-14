package config

import (
	gpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type PostgreSQL struct {
	DB *gorm.DB
}

func NewPostgreSQL(env string, cfg *PostgreSqlOption, dbLogger glogger.Interface) (*PostgreSQL, error) {
	logLevel := glogger.Warn
	if env == "local" {
		logLevel = glogger.Info
	}

	db, err := gorm.Open(gpostgres.Open(cfg.URI), &gorm.Config{
		Logger: dbLogger.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Pool)

	return &PostgreSQL{DB: db}, nil
}
