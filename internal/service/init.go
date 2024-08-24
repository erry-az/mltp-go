package service

import (
	"database/sql"
	"github.com/erry-az/mltp-go/db/query"
)

type Service struct {
	queries *query.Queries
	pool    *sql.DB
}

func New(queries *query.Queries, pool *sql.DB) *Service {
	return &Service{
		queries: queries,
		pool:    pool,
	}
}
