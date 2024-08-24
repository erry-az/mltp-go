package app

import (
	"database/sql"

	"github.com/erry-az/mltp-go/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func NewDBPool(config config.DbConfig) (*sql.DB, error) {
	db, err := otelsql.Open("pgx", config.Dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(config.MaxConnLifetime)
	db.SetConnMaxIdleTime(config.MaxConnIdleTime)
	db.SetMaxIdleConns(config.MinConn)
	db.SetMaxOpenConns(config.MaxConn)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
