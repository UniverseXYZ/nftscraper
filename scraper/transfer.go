package scraper

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/universexyz/nftscraper/contract/erc721"
)

func (s *Service) handleERC721Transfer(ctx context.Context, rawLog *ethlogscanner.Log) error {
	log := zerolog.Ctx(ctx)

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
		return errors.WithStack(err)
	}

	l.Raw = types.Log(*rawLog)

	log.Info().Str("from", l.From.String()).Str("to", l.To.String()).Str("token_id", l.TokenId.String()).Msg("erc721 transfer")

	return nil
}

func NewClient(projectId, projectSecret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	}
}

// authTransport decorates each request with a basic auth header.
type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func (t authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.ProjectId, t.ProjectSecret)
	return t.RoundTripper.RoundTrip(r)
}
