package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var files embed.FS

// New uses embedded SQL so migration paths do not depend on the working directory.
// The caller must close the returned migrator.
func New(databaseURL string) (*migrate.Migrate, error) {
	source, err := iofs.New(files, ".")
	if err != nil {
		return nil, fmt.Errorf("open embedded migrations: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		source.Close()
		return nil, fmt.Errorf("initialize migrations: %w", err)
	}
	return m, nil
}

// Up applies pending migrations before the service starts accepting work.
func Up(databaseURL string) (err error) {
	m, err := New(databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		err = errors.Join(err, sourceErr, databaseErr)
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
