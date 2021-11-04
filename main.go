package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/universexyz/nftscraper/contract/erc1155"
	"github.com/universexyz/nftscraper/contract/erc721"

	"github.com/universexyz/nftscraper/models" // how to make it local?
)

func init() {
	// load .env file in current working directory
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("unable to load .env file: `%s`", err.Error()))
	}

	// print error stack to the log messages
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
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

	// add logger to context
	ctx = logger.WithContext(ctx)

	// parse the given app flags
	flag.Parse()

	//db connection
	//refactor this approach with @Zero
	models.Init()	

	// execute app
	if err := run(ctx); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

// run is the entry point for the app, it should live in this function
func run(ctx context.Context) error {
	log := zerolog.Ctx(ctx)

	log.Info().Msg("hi")

	ethC, err := ethclient.DialContext(ctx, "https://mainnet.infura.io/v3/abc8f586485441c9b18cd4989f3951f8")
	if err != nil {
		return errors.WithStack(err)
	}

	cursor, err := loadCursor()
	if err != nil {
		return errors.WithStack(err)
	}

	erc721ABI, err := abi.JSON(strings.NewReader(erc721.ERC721ABI))
	if err != nil {
		return errors.WithStack(err)
	}

	erc1155ABI, err := abi.JSON(strings.NewReader(erc1155.ERC1155ABI))
	if err != nil {
		return errors.WithStack(err)
	}

	scanner, err := ethlogscanner.Scan(ctx, ethC,
		ethlogscanner.WithStayBehindToHead(8),
		ethlogscanner.WithStart(cursor),
		ethlogscanner.WithFilter(
			nil, [][]common.Hash{
				{
					erc721ABI.Events["Transfer"].ID,
					erc1155ABI.Events["TransferSingle"].ID,
					erc1155ABI.Events["TransferBatch"].ID,
				},
			}),
	)

	if err != nil {
		return errors.WithStack(err)
	}

	defer scanner.Close()

	for {
		select {
		case err := <-scanner.Done():
			return errors.WithStack(err)

		case err := <-scanner.Err():
			if errors.Is(err, context.Canceled) {
				return errors.WithStack(err)
			}

			log.Warn().Stack().Err(err).Msg("scan error")

		case l := <-scanner.Log():
			if l != nil {
				switch l.Topics[0] {
				case erc721ABI.Events["Transfer"].ID:
					log.Info().Msgf("Transfer: %s", l.Cursor().String())
					
					

					return errors.WithStack(nil)

				case erc1155ABI.Events["TransferSingle"].ID:
					log.Info().Msgf("TransferSingle: %s", l.Cursor().String())

				case erc1155ABI.Events["TransferBatch"].ID:
					log.Info().Msgf("TransferBatch: %s", l.Cursor().String())

				}
			}

		case n := <-scanner.Notify():
			if n != nil {
				switch n := n.(type) {
				case *ethlogscanner.ChunkSizeUpdated:
					log.Trace().Int("previous", n.Previous).Int("updated", n.Updated).Msg("chunk size updated")

				case *ethlogscanner.FilterStarted:
					log.Trace().Uint64("from", n.From).Uint64("to", n.To).Int("chunk", n.ChunkSize).Msg("filtering started")

				case *ethlogscanner.FilterCompleted:
					log.Debug().Uint64("from", n.From).Uint64("to", n.To).Int("chunk", n.ChunkSize).Dur("took", n.Duration).Bool("has_err", n.HasErr).Msg("filtering completed")

				case *ethlogscanner.CursorUpdated:
					//log.Trace().Uint64("block_num", n.Next.BlockNum()).Uint("tx_index", n.Next.TxIndex()).Uint("log_index", n.Next.LogIndex()).Msg("cursor updated")

					if err := storeCursor(n.Next); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
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
