package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Database struct {
	db *sqlx.DB
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

func New(ctx context.Context, dbPath string) (*Database, error) {
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if d == nil || d.db == nil {
		return nil
	}

	return d.db.Close()
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if d == nil || d.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	return d.db.ExecContext(ctx, query, args...)
}

func (d *Database) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	return d.db.GetContext(ctx, dest, query, args...)
}

func (d *Database) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	return d.db.SelectContext(ctx, dest, query, args...)
}

func (d *Database) WithTx(ctx context.Context, fn func(tx Executor) error) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction error: %w; rollback error: %v", err, rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
