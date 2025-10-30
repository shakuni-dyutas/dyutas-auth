package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionConfig struct {
	User string
	Pw   string
	Host string
	Port string
	DB   string
}

func NewRDBConnectionPool(cfg ConnectionConfig) (*Conn, error) {
	var connectionString = composeConnectionStringFrom(cfg)
	newPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	var newQueries = New(newPool)

	var newConn = &Conn{
		pool: newPool,
		Qrs:  newQueries,
	}

	return newConn, nil
}

func composeConnectionStringFrom(cfg ConnectionConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Pw, cfg.Host, cfg.Port, cfg.DB)
}

type Conn struct {
	pool *pgxpool.Pool
	Qrs  *Queries
}

type Trx struct {
	tx  pgx.Tx
	Qrs *Queries
}

func (c *Conn) BeginTx(ctx context.Context) (*Trx, error) {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	var trxQrs = c.Qrs.WithTx(tx)

	var newTrx = &Trx{tx: tx, Qrs: trxQrs}

	return newTrx, nil
}
