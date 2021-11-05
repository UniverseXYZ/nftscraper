package scraper

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/universexyz/nftscraper/contract/erc721"
)

func (s *Service) handleERC721Transfer(ctx context.Context, rawLog *ethlogscanner.Log) error {
	//	log := zerolog.Ctx(ctx)

	l := erc721.ERC721Transfer{}

	// fix non-indexed from, to, and token id
	if len(rawLog.Topics) == 1 && len(rawLog.Data) == common.HashLength*3 {
		rawLog.Topics = append(rawLog.Topics, common.BytesToHash(rawLog.Data))
		rawLog.Topics = append(rawLog.Topics, common.BytesToHash(rawLog.Data[common.HashLength:]))
		rawLog.Topics = append(rawLog.Topics, common.BytesToHash(rawLog.Data[common.HashLength*2:]))
		rawLog.Data = rawLog.Data[:0]
	}

	// fix non-indexed token id
	if len(rawLog.Topics) == 1+2 && len(rawLog.Data) == common.HashLength {
		rawLog.Topics = append(rawLog.Topics, common.BytesToHash(rawLog.Data))
		rawLog.Data = rawLog.Data[:0]
	}

	if err := s.erc721BC.UnpackLog(&l, s.evTransfer.RawName, types.Log(*rawLog)); err != nil {
		if err.Error() == "abi: attempting to copy no values while 3 arguments are expected" {

			fmt.Printf("EVENT ID%s\n", s.evTransfer.ID.String())
			for i, v := range rawLog.Topics {
				fmt.Printf("TOPIC[%d]: %s\n", i, v.String())
			}
			fmt.Println(common.Bytes2Hex(rawLog.Data))
		}
		return errors.WithStack(err)
	}

	//log.Info().Str("from", l.From.String()).Str("to", l.To.String()).Str("token_id", l.TokenId.String()).Msg("erc721 transfer")

	return nil
}
