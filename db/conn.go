package db

import (
	"context"
	"database/sql"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ctxKeyDB struct{}
type ctxKeyTX struct{}

const (
	defConnAppName       = "PIS"
	defConnPoolMin       = 2
	defConnPoolMax       = 50
	defConnMaxLifetime   = 30 * time.Minute
	defConnMaxIdleTime   = 3 * time.Minute
	defConnStmtCacheSize = 250
)

var ErrNoTx = errors.New("there is no active database transaction in the given context")

func Ctx(ctx context.Context) *sql.DB {
	if v, ok := ctx.Value(ctxKeyDB{}).(*sql.DB); v != nil && ok {
		return v
	}

	panic("here is a bug, there is no database instance in the given context")
}

func WithContext(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, ctxKeyDB{}, db)
}

// RunTx ...
func RunTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error {
	tx, ok := ctx.Value(ctxKeyTX{}).(*sql.Tx)
	if tx == nil || !ok {
		return errors.WithStack(ErrNoTx)
	}

	return errors.WithStack(f(ctx, tx))
}

// RunNewTx ...
func RunNewTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := Ctx(ctx).BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	defer tx.Rollback()

	if err := f(ctx, tx); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Connect ...
func Connect(ctx context.Context, dsn string) (*sql.DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	q := u.Query()

	confInt := func(key string, def int) (int, error) {
		if v, ok := q[key]; ok {
			if len(v) != 1 {
				return 0, errors.Errorf("multiple config value: `%s`", key)

			} else if v, err := strconv.Atoi(v[0]); err != nil {
				return 0, errors.WithStack(err)

			} else {
				return v, nil
			}
		}

		return def, nil
	}

	confDur := func(key string, def time.Duration) (time.Duration, error) {
		if v, ok := q[key]; ok {
			if len(v) != 1 {
				return 0, errors.Errorf("multiple config value: `%s`", key)

			} else if v, err := time.ParseDuration(v[0]); err != nil {
				return 0, errors.WithStack(err)

			} else {
				return v, nil
			}
		}

		return def, nil
	}

	minConn, err := confInt("pool_min_conns", defConnPoolMin)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	maxConn, err := confInt("pool_max_conns", defConnPoolMax)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	maxLife, err := confDur("pool_max_conn_lifetime", defConnMaxLifetime)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	maxIdle, err := confDur("pool_max_conn_idle_time", defConnMaxIdleTime)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if _, ok := q["application_name"]; !ok {
		q.Add("application_name", defConnAppName)
	}

	if _, ok := q["statement_cache_capacity"]; !ok {
		q.Add("statement_cache_capacity", strconv.Itoa(defConnStmtCacheSize))
	}

	u.RawQuery = q.Encode()

	dbConf, err := pgx.ParseConfig(u.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dbConf.Logger = zerologadapter.NewLogger(zerolog.Ctx(ctx).With().Str("comp", "db").Logger())

	db := stdlib.OpenDB(*dbConf)

	db.SetMaxIdleConns(minConn)
	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxIdleTime(maxIdle)
	db.SetConnMaxLifetime(maxLife)

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
