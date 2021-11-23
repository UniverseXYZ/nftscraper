package store

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/universexyz/nftscraper/db"
	"github.com/universexyz/nftscraper/model"
)

type TransferStore interface {
	// FindByTokenID(ctx context.Context, contractAddr, tokenID string) (*model.NFT, error)
	Save(ctx context.Context, transfer *model.Transfer) error
}

type transferStore struct {
	stmtSave *sql.Stmt
}

// Creates and return an instance of transferStore which implements TransferStore interface.
func NewTransferStore(ctx context.Context) (TransferStore, error) {
	var err error

	dbConn := db.Ctx(ctx)

	store := &transferStore{}

	store.stmtSave, err = dbConn.PrepareContext(ctx, `
		INSERT INTO transfer (
			id,
			contract_addr,
			token_id,
			"from",
			"to",
			amount,
			"type",
			"tx_hash",
			"log_index"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return store, nil
}

// Adds a new entry to the nft table
func (t *transferStore) Save(ctx context.Context, transfer *model.Transfer) error {
	return db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.StmtContext(ctx, t.stmtSave).ExecContext(ctx, 
			transfer.ID,
			transfer.ContractAddress,
			transfer.TokenID,
			transfer.From,
			transfer.To,
			transfer.Amount,
			transfer.Type,
			transfer.TxHash,
			transfer.LogIndex)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}
