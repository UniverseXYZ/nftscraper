package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/universexyz/nftscraper/models"
)

type Store interface {
	AddNFT(context.Context, *models.NFT) (uuid.UUID, error)
	AddTransfer(context.Context, *models.NFT) (uuid.UUID, error)
	AddNFTCollection(context.Context, *models.NFTCollection) (uuid.UUID, error)
}
