package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

func New(ctx context.Context, databaseURL string) (*Database, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := pingWithRetry(ctx, db, 30*time.Second, time.Second); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Database{db: db}, nil
}

func pingWithRetry(ctx context.Context, db *sqlx.DB, timeout, interval time.Duration) error {
	retryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastErr error
	for {
		if err := db.PingContext(retryCtx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-retryCtx.Done():
			return fmt.Errorf("database did not become ready: %w", lastErr)
		case <-time.After(interval):
		}
	}
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
