package scraper

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/universexyz/nftscraper/conf"
	"github.com/universexyz/nftscraper/contract/erc1155"
	"github.com/universexyz/nftscraper/contract/erc721"
	"github.com/universexyz/nftscraper/store"
)

type ServiceDeps interface {
	ScraperStore() store.ScraperStore
}

type Service struct {
	deps             ServiceDeps
	conf             conf.Config
	ethC             *ethclient.Client
	run              int64
	erc721ABI        *abi.ABI
	erc1155ABI       *abi.ABI
	evTransfer       abi.Event
	evTransferSingle abi.Event
	evTransferBatch  abi.Event
	erc721BC         *bind.BoundContract
	erc1155BC        *bind.BoundContract
}

const scraperName = "ethereum:mainnet"

func NewService(ctx context.Context, deps ServiceDeps) (*Service, error) {
	s := &Service{
		deps: deps,
		conf: conf.Conf(),
	}

	ethC, err := ethclient.DialContext(ctx, s.conf.Web3URL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.ethC = ethC

	erc721ABI, err := abi.JSON(strings.NewReader(erc721.ERC721ABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.erc721ABI = &erc721ABI

	erc1155ABI, err := abi.JSON(strings.NewReader(erc1155.ERC1155ABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.erc1155ABI = &erc1155ABI

	s.evTransfer = s.erc721ABI.Events["Transfer"]
	s.evTransferSingle = s.erc1155ABI.Events["TransferSingle"]
	s.evTransferBatch = s.erc1155ABI.Events["TransferBatch"]

	s.erc721BC = bind.NewBoundContract(common.Address{}, *s.erc721ABI, nil, nil, nil)
	s.erc1155BC = bind.NewBoundContract(common.Address{}, *s.erc1155ABI, nil, nil, nil)

	return s, nil
}

func (s *Service) Run(ctx context.Context, shutdownCh <-chan struct{}) error {
	if !atomic.CompareAndSwapInt64(&s.run, 0, 1) {
		return errors.New("service is already running")
	}

	defer atomic.StoreInt64(&s.run, 0)

	ticker := time.NewTicker(s.conf.ChainScanPeriod)

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := zerolog.Ctx(ctx)

	for {
		err := s.scan(subCtx, shutdownCh)
		if err != nil {
			log.Warn().Stack().Err(err).Msg("an error occurred while scanning the chain data")
		}

		select {
		case <-ctx.Done():
			return errors.WithStack(ctx.Err())

		case <-shutdownCh:
			return nil

		case _, ok := <-ticker.C:
			if !ok {
				return nil
			}

			continue
		}
	}
}

func (s *Service) scan(ctx context.Context, shutdownCh <-chan struct{}) error {
	ctx, cancel := context.WithTimeout(ctx, s.conf.ChainScanPeriod-(10*time.Second))
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-shutdownCh:
			cancel()
		}
	}()

	log := zerolog.Ctx(ctx)

	cursor, err := s.deps.ScraperStore().LoadCursor(ctx, scraperName)
	if err != nil {
		return errors.WithStack(err)
	}

	scanner, err := ethlogscanner.Scan(ctx, s.ethC,
		ethlogscanner.WithStayBehindToHead(int(s.conf.StayBehindToHead)),
		ethlogscanner.WithStart(cursor),
		ethlogscanner.WithFilter(
			nil, [][]common.Hash{
				{
					s.evTransfer.ID,
					s.evTransferSingle.ID,
					s.evTransferBatch.ID,
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
			if errors.Is(err, context.DeadlineExceeded) {
				return nil
			}

			return errors.WithStack(err)

		case err := <-scanner.Err():
			if err != nil {
				if errors.Is(err, context.Canceled) {
					select {
					case <-shutdownCh:
						return nil

					default:
						return errors.WithStack(err)
					}
				}

				if errors.Is(err, context.DeadlineExceeded) {
					return nil
				}

				log.Warn().Stack().Err(err).Msg("scan error")
			}

		case l := <-scanner.Log():
			if l != nil {
				logMsg :=
					log.Trace().
						Uint64("block", l.Cursor().BlockNum()).
						Uint("tx_index", l.Cursor().TxIndex()).
						Uint("log_index", l.Cursor().LogIndex())

				switch l.Topics[0] {
				case s.evTransfer.ID:
					logMsg.Msg("found an erc721 transfer event")

					err := s.handleERC721Transfer(ctx, l)
					if err != nil {
						return errors.WithStack(err)
					}

				case s.evTransferSingle.ID:
					logMsg.Msg("found an erc1155 single transfer event")

					err := s.handleERC1155TransferSingle(ctx, l)
					if err != nil {
						return errors.WithStack(err)
					}

				case s.evTransferBatch.ID:
					logMsg.Msg("found an erc1155 batch transfer event")

					err := s.handleERC1155TransferBatch(ctx, l)
					if err != nil {
						return errors.WithStack(err)
					}

				default:
					logMsg.Msgf("found an unexpected event and skipped. event id: %s", l.Topics[0].String())
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
					log.Trace().Uint64("block_num", n.Next.BlockNum()).Uint("tx_index", n.Next.TxIndex()).Uint("log_index", n.Next.LogIndex()).Msg("cursor updated")

					if err := s.deps.ScraperStore().StoreCursor(ctx, scraperName, n.Next); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}
}
