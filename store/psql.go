package store

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/universexyz/nftscraper/conf"
	"github.com/universexyz/nftscraper/constants"
	"github.com/universexyz/nftscraper/models"

	"github.com/jmoiron/sqlx"
)

type client struct {
	DB *sqlx.DB
}

// Creates new instance of the DB connection.
func NewClient() *client {
	connection, err := sqlx.Open(constants.POSTGRES, conf.Conf().PostgresDSN)

	if err != nil {
		log.Fatalf("Cannot connect to DB " + err.Error())
	}

	return &client{DB: connection}
}

// Adds an entry to the transfer table
func (c *client) AddTransfer(ctx context.Context, transfer models.Transfer) (uuid.UUID, error) {
	var newRowID uuid.UUID = uuid.Nil
	err := c.DB.QueryRowContext(ctx, `
			INSERT INTO transfer (
				contractAddress,
				tokenId,
				"from",
				"to",
				amount,
				"type",
				txHash,
				logIndex
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`,
		transfer.ContractAddress,
		transfer.TokenID,
		transfer.From,
		transfer.To,
		transfer.Amount,
		transfer.Type,
		transfer.TxHash,
		transfer.LogIndex).Scan(&newRowID)

	if(err != nil) {
		return uuid.Nil, err
	}

	return newRowID, nil
}

// Adds an entry to the nft table
func (c *client) AddNFT(ctx context.Context, NFT models.NFT) (uuid.UUID, error) {
	var newRowID uuid.UUID = uuid.Nil
	err := c.DB.QueryRowContext(ctx, `
			INSERT INTO nft (
				contractAddress, 
				tokenId, 
				ownerAddress, 
				name, 
				symbol, 
				tokenUri,
				optimizedUrl, 
				thumbnailUrl, 
				attributes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`,
		NFT.ContractAddress,
		NFT.TokenID,
		NFT.OwnerAddress,
		NFT.Name,
		NFT.Symbol,
		NFT.TokenURI,
		NFT.OptimizedURL,
		NFT.ThumbnailURL,
		NFT.Attributes).Scan(&newRowID)
	
	if err != nil {
		return uuid.Nil, err
	}
	
	return newRowID, nil
}

// Adds an entry to the nftCollection table
func (c *client) AddNFTCollection(ctx context.Context, NFT models.NFTCollection) (uuid.UUID, error) {
	var newRowID uuid.UUID = uuid.Nil
	err := c.DB.QueryRowContext(ctx, `
			INSERT INTO nft (
				contractAddress, 
				name, 
				symbol, 
				numberOfNfts
			)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`,
		NFT.ContractAddress,
		NFT.Name,
		NFT.Symbol,
		NFT.NumberOfNFTs).Scan(&newRowID)
	
	if err != nil {
		return uuid.Nil, err
	}
	
	return newRowID, nil
}

