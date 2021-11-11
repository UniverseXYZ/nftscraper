package store

import (
	"context"
	"database/sql"

	"github.com/mgurevin/ethlogscanner"
	"github.com/pkg/errors"
	"github.com/universexyz/nftscraper/db"
)

type scraperStore struct {
	stmtStoreCursor *sql.Stmt
	stmtLoadCursor  *sql.Stmt
}

func NewScraperStore(ctx context.Context) (ScraperStore, error) {
	var err error

	dbConn := db.Ctx(ctx)

	store := &scraperStore{}

	store.stmtLoadCursor, err = dbConn.PrepareContext(ctx, "select last_block_num, last_tx_num, last_log_index from scraper_cursor where name=$1")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	store.stmtStoreCursor, err = dbConn.PrepareContext(ctx, "insert into scraper_cursor(name, last_block_num, last_tx_num, last_log_index) values($1, $2, $3, $4) on conflict(name) do update set last_block_num=EXCLUDED.last_block_num, last_tx_num=EXCLUDED.last_tx_num, last_log_index=EXCLUDED.last_log_index")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return store, nil
}

func (s *scraperStore) LoadCursor(ctx context.Context, scraperName string) (ethlogscanner.Cursor, error) {
	var lastBlockNum, lastTxIndex, lastLogIndex int

	err := db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.StmtContext(ctx, s.stmtLoadCursor).QueryRowContext(ctx, scraperName)
		if err := row.Scan(&lastBlockNum, &lastTxIndex, &lastLogIndex); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}

			return errors.WithStack(err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return ethlogscanner.MakeCursor(uint64(lastBlockNum), uint(lastTxIndex), uint(lastLogIndex)), nil
}

func (s *scraperStore) StoreCursor(ctx context.Context, scraperName string, cursor ethlogscanner.Cursor) error {
	return db.RunTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.StmtContext(ctx, s.stmtStoreCursor).ExecContext(ctx, scraperName, cursor.BlockNum(), cursor.TxIndex(), cursor.LogIndex())
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}
