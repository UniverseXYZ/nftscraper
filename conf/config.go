package conf

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type logLevel zerolog.Level

type Config struct {
	PostgresDSN      string        `env:"POSTGRESQL_DSN"`
	DBHost      	 string        `env:"DB_HOST"`
	DBPort      	 uint32        `env:"DB_PORT"`
	DBUser      	 string        `env:"DB_USER"`
	DBPassword       string        `env:"DB_PASSWORD"`
	DBName      	 string        `env:"DB_NAME"`
	Web3URL          string        `env:"WEB3_URL"`
	ChainScanPeriod  time.Duration `env:"CHAIN_SCAN_PERIOD"   envDefault:"30s"`
	StayBehindToHead uint          `env:"STAY_BEHIND_TO_HEAD" envDefault:"8"`
	LogLevel         logLevel      `env:"LOG_LEVEL"           envDefault:"DEBUG"`
}

var once sync.Once
var config Config

func Parse() (err error) {
	// try to parse env vars only once
	once.Do(func() {
		/// load the .env file in the current working directory if exists
		if e := godotenv.Load(); e != nil && !errors.Is(e, os.ErrNotExist) {
			err = errors.Wrapf(err, "unable to load .env file: `%s`", e.Error())
			return
		}

		// parse env variables
		err = errors.WithStack(env.Parse(&config))
	})

	return
}

func Conf() Config {
	if err := Parse(); err != nil {
		panic(fmt.Sprintf("unable to parse environment variables: %+v", err))
	}

	return config
}

func (l *logLevel) UnmarshalText(text []byte) error {
	lvl, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(string(text))))
	if err != nil {
		return errors.Wrapf(err, "unable to parse log level `%s`: %s", string(text), err.Error())
	}

	*l = logLevel(lvl)

	return nil
}
