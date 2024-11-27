package repository

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type DataStore interface {
	Atomic(ctx context.Context, fn func(DataStore) error) error
}

type dataStore struct {
	conn *sql.DB
	db   DBTX
}

func NewDataStore(db *sql.DB) DataStore {
	return &dataStore{
		conn: db,
		db:   db,
	}
}

func (s *dataStore) Atomic(ctx context.Context, fn func(DataStore) error) error {
	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(&dataStore{conn: s.conn, db: tx})
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
