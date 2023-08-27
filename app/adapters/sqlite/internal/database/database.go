package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pietjan/tm/app/adapters/sqlite/internal/driver"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
)

func New(dsn string) Database {
	db, err := sql.Open(`sqlite`, dsn)
	if err != nil {
		panic(err)
	}

	return Database{
		db: db,
		qb: driver.New(db),
	}
}

type Database struct {
	db *sql.DB
	qb qbdb.DB
}

func (d Database) DB() *sql.DB {
	return d.db
}

func (d Database) QB() qbdb.DB {
	return d.qb
}

func (d Database) ExecContext(ctx context.Context, q qb.Query) (res sql.Result, err error) {
	query, args := d.qb.Render(q)
	defer func() {
		if err != nil {
			fmt.Println(query)
		}
	}()

	res, err = d.db.ExecContext(ctx, query, args...)

	return
}

func (d Database) QueryContext(ctx context.Context, q qb.Query) (rows *sql.Rows, err error) {
	query, args := d.qb.Render(q)
	defer func() {
		if err != nil {
			fmt.Println(query)
		}
	}()

	rows, err = d.db.QueryContext(ctx, query, args...)
	return
}

func (d Database) QueryRowContext(ctx context.Context, q qb.Query) *sql.Row {
	query, args := d.qb.Render(q)

	return d.db.QueryRowContext(ctx, query, args...)
}

type Tx struct {
	qb qbdb.DB
	tx *sql.Tx
}

func (t Tx) ExecContext(ctx context.Context, q qb.Query) (res sql.Result, err error) {
	query, args := t.qb.Render(q)
	res, err = t.tx.ExecContext(ctx, query, args...)

	return
}

func (t Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t Tx) Commit() error {
	return t.tx.Commit()
}

func (d Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	t, err := d.db.BeginTx(ctx, opts)
	return Tx{
		qb: d.qb,
		tx: t,
	}, err
}
