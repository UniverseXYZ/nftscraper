package scraper

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/universexyz/nftscraper/contract/erc1155"
)

func (s *Service) handleERC1155TransferBatch(ctx context.Context, rawLog *ethlogscanner.Log) error {
	log := zerolog.Ctx(ctx)

	l := erc1155.ERC1155TransferBatch{}

	if err := s.erc1155BC.UnpackLog(&l, s.evTransferBatch.RawName, types.Log(*rawLog)); err != nil {
		return errors.WithStack(err)
	}

	ids := make([]string, len(l.Ids))
	vals := make([]string, len(l.Values))

	for i, v := range l.Ids {
		ids[i] = v.String()
	}

	for i, v := range l.Values {
		vals[i] = v.String()
	}

	log.Info().
		Str("operator", l.Operator.String()).
		Str("from", l.From.String()).
		Str("to", l.To.String()).
		Strs("ids", ids).
		Strs("values", vals).
		Msg("erc1155 batch transfer")

	return nil
}
