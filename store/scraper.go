package store

import (
	"context"

	"github.com/mgurevin/ethlogscanner"
)

type ScraperStore interface {
	LoadCursor(ctx context.Context, scraperName string) (ethlogscanner.Cursor, error)
	StoreCursor(ctx context.Context, scraperName string, cursor ethlogscanner.Cursor) error
}
