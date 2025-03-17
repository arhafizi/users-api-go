package dbConf

import (
	"database/sql"
	"fmt"

	"example.com/api/config"
	"example.com/api/pkg/logging"
)

func InitDb(cfg *config.Config, logger logging.ILogger) *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal(logging.Postgres, logging.SubCategory(logging.Connection), "Failed to connect to database", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
		})
		return nil
	}

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		logger.Fatal(logging.Postgres, logging.SubCategory(logging.Connection), "Failed to ping database", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
		})
	}
	logger.Info(logging.Postgres, logging.SubCategory(logging.Internal), "Successfully connected to database", nil)

	return db
}
