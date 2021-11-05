package migrate

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/universexyz/nftscraper/conf"
)

//Executes DB migration
func Run(migrationType string) {
	m, err := migrate.New(
		"file://migrate/migrations/",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			conf.Conf().DBUser,
			conf.Conf().DBPassword,
			conf.Conf().DBHost,
			conf.Conf().DBPort,
			conf.Conf().DBName,
			conf.Conf().DBSSL))
	if err != nil {
		log.Fatal("Error trying to prepare migration: " + err.Error())
	}

	switch migrationType {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal("Error while running migration up: " + err.Error())
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal("Error while running migration down: " + err.Error())
		}
	default:
		log.Fatal("Usage: go run main migrate up/down")
	}
}
