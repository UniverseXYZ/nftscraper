package store

import (
	"context"

	"github.com/universexyz/nftscraper/model"
)

type NFTStore interface {
	FindByTokenID(ctx context.Context, contractAddr, tokenID string) (*model.NFT, error)
	Save(ctx context.Context, nft *model.NFT) error
}
