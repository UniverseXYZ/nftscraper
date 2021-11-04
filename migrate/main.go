package migrate

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/universexyz/nftscraper/migrate" // @Zero how to make it local?
	"github.com/universexyz/nftscraper/models" // @Zero how to make it local?
)

func main() {

	models.Init();

	sqlMigration := migrate.Sqlx{
		Migrations: []migrate.SqlxMigration{
			migrate.SqlxFileMigration("migrate_01", "migrations/migrate_01.sql", "migrations/migrate_01.undo.sql"),
			migrate.SqlxMigration{
			ID: "002_add_currency",
			Migrate: func(tx *sqlx.Tx) error {
					// add currency field, then fill existing entries with "USD"
					return nil
				},
			Rollback: func(tx *sqlx.Tx) error {
					// drop the currency field
					return nil
				},
			},
		},
	}

	migrate.Migrate(models.DB, "postgres");
	  
}