package store

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/universexyz/nftscraper/db"
	"github.com/universexyz/nftscraper/model"
)

type NFTCollectionStore interface {
	Save(ctx context.Context, nft *model.NFTCollection) error
	FindByContractAddress(ctx context.Context, contractAddress string) ([]model.NFTCollection, error)
}

type nftCollectionStore struct {
	stmtSave *sql.Stmt
	stmtFindByContractAddress *sql.Stmt
}

// Creates and return an instance of nftCollectionStore which implements NFTCollectionStore interface.
func NewNFTCollectionStore(ctx context.Context) (NFTCollectionStore, error) {
	var err error

	dbConn := db.Ctx(ctx)

	store := &nftCollectionStore{}

	store.stmtSave, err = dbConn.PrepareContext(ctx, `
		INSERT INTO nft_collection (
			id,
			contract_addr,
			name,
			symbol,
			num_nfts
		)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	store.stmtFindByContractAddress, err = dbConn.PrepareContext(ctx, `
		SELECT *
		FROM nft_collection
		WHERE contract_addr = $1
	`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return store, nil
}

// Adds a new entry to the nft_collection table
func (n *nftCollectionStore) Save(ctx context.Context, nftCollection *model.NFTCollection) error {
	return db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.StmtContext(ctx, n.stmtSave).ExecContext(ctx,
			nftCollection.ID,
			nftCollection.ContractAddress,
			nftCollection.Name,
			nftCollection.Symbol,
			nftCollection.NumberOfNFTs)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

// Returns a slice of model.NFTCollection by looking up in the nft_collection table by contractAddress
// There may be multiple nft collections with same contract address due to artblocks.
// @See https://www.artblocks.io/
func (n *nftCollectionStore) FindByContractAddress(ctx context.Context, contractAddress string) ([]model.NFTCollection, error) {
	var nftCollections = make([]model.NFTCollection, 0)

	err := db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		rows, err := tx.StmtContext(ctx, n.stmtFindByContractAddress).QueryContext(ctx, contractAddress)
		if err != nil {
			return errors.WithStack(err)
		}

		defer rows.Close()

		for rows.Next() {
			var nftCollection model.NFTCollection
			err := rows.Scan(
				&nftCollection.ID,
				&nftCollection.ContractAddress,
				&nftCollection.Name,
				&nftCollection.Symbol,
				&nftCollection.NumberOfNFTs)
			if err != nil {
				// return errors.WithStack(err)
				break
			}
			nftCollections = append(nftCollections, nftCollection)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nftCollections, nil
}