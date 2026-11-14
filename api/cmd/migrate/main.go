package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BoaglioFranco/casita/api/migrations"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	action := flag.String("action", "up", "migration action: up, down, force")
	version := flag.Int("version", -1, "verified schema version to record with force")
	flag.Parse()

	if *action != "up" && *action != "down" && *action != "force" {
		return fmt.Errorf("unknown migration action %q", *action)
	}
	if *action == "force" && *version < 0 {
		return errors.New("force requires a nonnegative -version")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	m, err := migrations.New(databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		err = errors.Join(err, sourceErr, databaseErr)
	}()

	switch *action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "force":
		err = m.Force(*version)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
