package store

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/universexyz/nftscraper/db"
	"github.com/universexyz/nftscraper/model"
)

type NFTStore interface {
	Save(ctx context.Context, nft *model.NFT) error
	FindByContractAddressAndTokenID(ctx context.Context, contractAddress string, tokenID string) (*model.NFT, error)
}

type nftStore struct {
	stmtSave *sql.Stmt
	stmtFindByContractAddressAndTokenID *sql.Stmt
}

// Creates and return an instance of nftStore which implements NFTStore interface.
func NewNFTStore(ctx context.Context) (NFTStore, error) {
	var err error

	dbConn := db.Ctx(ctx)

	store := &nftStore{}

	store.stmtSave, err = dbConn.PrepareContext(ctx, `
		INSERT INTO nft (
			id,
			nft_collection_id,
			contract_addr,
			token_id,
			owner_addr,
			"name",
			symbol,
			token_uri,
			optimized_url,
			thumbnail_url,
			attributes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	store.stmtFindByContractAddressAndTokenID, err = dbConn.PrepareContext(ctx, `
		SELECT *
		FROM nft
		WHERE contract_addr = $1 AND
			token_id = $2
	`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return store, nil
}

// Adds a new entry to the nft table
func (n *nftStore) Save(ctx context.Context, nft *model.NFT) error {
	return db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.StmtContext(ctx, n.stmtSave).ExecContext(ctx,
			nft.ID,
			nft.NFTCollectionID,
			nft.ContractAddress,
			nft.TokenID,
			nft.OwnerAddress,
			nft.Name,
			nft.Symbol,
			nft.TokenURI,
			nft.OptimizedURL,
			nft.ThumbnailURL,
			nft.Attributes)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

// Returns an instance of model.NFT by looking up in the nft table by contractAddress and tokenID
func (n *nftStore) FindByContractAddressAndTokenID(ctx context.Context, contractAddress string, tokenID string) (*model.NFT, error) {
	var nft model.NFT
	
	err := db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.StmtContext(ctx, n.stmtFindByContractAddressAndTokenID).QueryRowContext(ctx, contractAddress, tokenID)
		err := row.Scan(
			&nft.ID,
			&nft.NFTCollectionID,
			&nft.ContractAddress,
			&nft.TokenID,
			&nft.OwnerAddress,
			&nft.Name,
			&nft.Symbol,
			&nft.TokenURI,
			&nft.OptimizedURL,
			&nft.ThumbnailURL,
			&nft.Attributes) 
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}

			return errors.WithStack(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &nft, nil
}
