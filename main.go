package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/universexyz/nftscraper/conf"
	"github.com/universexyz/nftscraper/constants"
	"github.com/universexyz/nftscraper/migrate"
	"github.com/universexyz/nftscraper/scraper"
	"github.com/universexyz/nftscraper/migrate"
)

var migrationType string

func init() {
	// print error stack to the log messages
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	flag.StringVar(&migrationType, constants.MIGRATE, "", "-migrate up or -migrate down")
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

	// If there's "migrate" argument then we only run DB migration
	if(len(os.Args) > 1 && "migrate" == os.Args[1]) {
		migrate.Run()
		os.Exit(0);
	}

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

	// If there's "migrate" argument then we only run DB migration
	if len(migrationType) > 0 {
		err := migrate.Run(ctx, migrationType)
		if(err != nil) {
			return errors.WithStack(err)
		}
		os.Exit(0)
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
