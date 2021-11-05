package store

import (
	"log"

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
func (c *client) AddTransfer(transfer models.Transfer) {
	
}

// @Zero - can't understand why this is visible when there's no AddNft declared in the interface
// Adds an entry to the nft table
func (c *client) AddNft(nft models.Nft) (string, string) {
	contractAddress := ""
	tokenId := ""
	err := c.DB.QueryRow(`
			INSERT INTO nft (
				contractAddress, 
				tokenId, 
				ownerAddress, 
				name, 
				symbol, 
				optimizedUrl, 
				thumbnailUrl, 
				attributes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING contractAddress, tokenId
		`,
		nft.ContractAddress,
		nft.TokenId,
		nft.OwnerAddress,
		nft.Name,
		nft.Symbol,
		nft.OptimizedUrl,
		nft.ThumbnailUrl,
		nft.Attributes).Scan(&contractAddress, &tokenId)
	
	if err != nil {
		panic(err)
	}
	
	return contractAddress, tokenId
}

