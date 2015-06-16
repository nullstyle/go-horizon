package db

import (
	"github.com/jmoiron/sqlx"
	sq "github.com/lann/squirrel"
	"github.com/stellar/go-horizon/log"
	"golang.org/x/net/context"
)

type SqlQuery struct {
	DB *sqlx.DB
}

func (q SqlQuery) Select(ctx context.Context, sql sq.SelectBuilder, dest interface{}) error {
	sql = sql.PlaceholderFormat(sq.Dollar)
	query, args, err := sql.ToSql()

	if err != nil {
		return err
	}

	log.WithField(ctx, "sql", query).Info("Executing query")

	return q.DB.Select(dest, query, args...)
}

func (q SqlQuery) Get(ctx context.Context, sql sq.SelectBuilder, dest interface{}) error {
	sql = sql.PlaceholderFormat(sq.Dollar)
	query, args, err := sql.ToSql()

	if err != nil {
		return err
	}
	return q.DB.Get(dest, query, args...)
}
