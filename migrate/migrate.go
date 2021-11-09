package migrate

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/universexyz/nftscraper/conf"
)

//Executes DB migration
func Run(ctx context.Context, migrationType string) error {
	m, err := migrate.New(
		"file://migrate/migrations/",
		conf.Conf().PostgresDSN)
	if err != nil {
		return err
	}

	switch migrationType {
	case "up":
		if err := m.Up(); err != nil {
			return err
		}
	case "down":
		if err := m.Down(); err != nil {
			return err
		}
	default:
		return errors.New("Usage: go run main.go -migrate up or go run main.go -migrate down")
	}

	return nil
}
