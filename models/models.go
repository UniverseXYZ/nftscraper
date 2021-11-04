package models

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/universexyz/nftscraper/config" // @Zero how to make it local?

	_ "github.com/lib/pq" // @Zero this package looks outdated
)

var DB *sqlx.DB

func Init() *sqlx.DB {
	DB, err := sqlx.Open("postgres", "postgres://" + config.Get("DB_USER") + ":" + config.Get("DB_PASSWORD") + "@" + config.Get("DB_HOST"))

	if err != nil {
		log.Fatalf("Cannot connect to DB " + err.Error())
	}

	return DB
}

type Transfer struct {
	ContractAddress string `json:"contractAddress"`
	TokenId         uint64 `json:"tokenId"` //should it be string?
	From            string `json:"from"`
	To              string `json:"to"`
	Amount          uint64 `json:"amount"`
	Type            string `json:"type"`
	TxHash 			string `json:"txHash"`
	LogIndex        uint64 `json:"logIndex"`
}
