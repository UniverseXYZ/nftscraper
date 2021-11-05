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
func Run() {

	m, err := migrate.New(
		"file://migrate/migrations",
		// conf.Conf().PostgresDSN)
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Conf().DBHost, 
		conf.Conf().DBPort, 
		conf.Conf().DBUser, 
		conf.Conf().DBPassword, 
		conf.Conf().DBName))
	if err != nil {
		log.Fatal("Error trying to start migration: " + err.Error())
	}

	if err := m.Up(); err != nil {
		log.Fatal("Error while running migration: " + err.Error())
	}
}