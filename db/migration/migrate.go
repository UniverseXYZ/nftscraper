package migration

import (
	"context"
	"database/sql"
	"embed"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	migrationSchemaName  = "public"
	migrationTableName   = "_db_migrations"
	migrationStmtTimeout = 5 * time.Minute
)

//go:embed migrations
var migrationFS embed.FS

type migrationSource struct {
	httpfs.PartialDriver
}

type migrationLogger struct {
	l zerolog.Logger
}

// NewMigrate ...
func NewMigrate(ctx context.Context, db *sql.DB) (*migrate.Migrate, error) {
	sourceDriver := &migrationSource{}

	if err := sourceDriver.Init(http.FS(migrationFS), "migrations"); err != nil {
		return nil, errors.WithStack(err)
	}

	targetDriver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable:  migrationTableName,
		SchemaName:       migrationSchemaName,
		StatementTimeout: migrationStmtTimeout,
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	m, err := migrate.NewWithInstance("migrations", sourceDriver, "database", targetDriver)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	m.Log = &migrationLogger{l: zerolog.Ctx(ctx).With().Str("comp", "db-migration").Logger()}

	return m, nil
}

func (f *migrationSource) Open(filename string) (source.Driver, error) {
	src := &migrationSource{}

	if err := src.Init(http.FS(migrationFS), filename); err != nil {
		return nil, errors.WithStack(err)
	}

	return src, nil
}

func (l *migrationLogger) Printf(format string, v ...interface{}) {
	l.l.Info().Msgf(format, v...)
}

func (l *migrationLogger) Verbose() bool {
	return true
}
