package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"

	_ "github.com/joho/godotenv"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/universexyz/nftscraper/conf"
	"github.com/universexyz/nftscraper/db"
	"github.com/universexyz/nftscraper/db/migration"
	"github.com/universexyz/nftscraper/scraper"
)

var argMigrate string

func init() {
	// print error stack to the log messages
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	flag.StringVar(&argMigrate, "migrate", "", "start database migration UP or DOWN")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// send all log messages to stdout
	logWriter := zerolog.SyncWriter(os.Stdout)

	// if the program starts on the console use styled log writer
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logWriter = zerolog.NewConsoleWriter()
	}

	// init logger
	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	if err := conf.Parse(); err != nil {
		logger.Fatal().Stack().Err(err).Msg("configuration error")
	}

	// set the minimum log severity to output
	logger = logger.Level(zerolog.Level(conf.Conf().LogLevel))

	// add logger to context
	ctx = logger.WithContext(ctx)

	// parse the given app flags
	flag.Parse()

	// execute app
	if err := run(ctx); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

type x struct{}

func (xx *x) LoadChainScannerCursor(ctx context.Context) (ethlogscanner.Cursor, error) {
	return loadCursor()
}
func (xx *x) StoreChainScannerCursor(ctx context.Context, cursor ethlogscanner.Cursor) error {
	return storeCursor(cursor)
}

// run is the entry point for the app, it should live in this function
func run(ctx context.Context) error {
	dbConn, err := db.Connect(ctx, conf.Conf().PostgresDSN)
	if err != nil {
		return errors.WithStack(err)
	}

	defer dbConn.Close()

	// add database connection to the application context
	ctx = db.WithContext(ctx, dbConn)

	// execute migrations if neeeded
	if err := startMigration(ctx); err != nil {
		return errors.WithStack(err)
	}

	s, err := scraper.NewService(ctx, &x{})
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.Run(ctx, make(<-chan struct{}))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func loadCursor() (ethlogscanner.Cursor, error) {
	if _, err := os.Stat("cursor.json"); os.IsNotExist(err) {
		return 0, nil
	}

	f, err := os.OpenFile("cursor.json", os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	var c ethlogscanner.Cursor

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return 0, err
	}

	return c, nil
}

func storeCursor(c ethlogscanner.Cursor) error {
	f, err := os.OpenFile("cursor.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(c); err != nil {
		return err
	}

	return nil
}

// startMigration executes the migration process if requested by the user
func startMigration(ctx context.Context) error {
	doMigrate := false

	// check if the flag is provided
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "migrate" {
			doMigrate = true
		}
	})

	if !doMigrate {
		return nil
	}

	migrate, err := migration.NewMigrate(ctx, db.Ctx(ctx))
	if err != nil {
		return errors.WithStack(err)
	}

	switch strings.ToLower(strings.TrimSpace(argMigrate)) {
	case "up":
		if err = migrate.Up(); err != nil {
			return errors.WithStack(err)
		}

	case "down":
		if err = migrate.Down(); err != nil {
			return errors.WithStack(err)
		}

	default:
		return errors.Errorf("unable to parse direction of the migration: `%s` - only `up` or `down` supported", argMigrate)
	}

	return nil
}
