package scraper

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/universexyz/nftscraper/contract/erc1155"
)

func (s *Service) handleERC1155TransferSingle(ctx context.Context, rawLog *ethlogscanner.Log) error {
	log := zerolog.Ctx(ctx)

	l := erc1155.ERC1155TransferSingle{}

	if err := s.erc1155BC.UnpackLog(&l, s.evTransferSingle.RawName, types.Log(*rawLog)); err != nil {
		return errors.WithStack(err)
	}

	log.Info().
		Str("operator", l.Operator.String()).
		Str("from", l.From.String()).
		Str("to", l.To.String()).
		Str("id", l.Id.String()).
		Str("value", l.Value.String()).
		Msg("erc1155 single transfer")

	return nil
}
