package model

import "github.com/google/uuid"

type Transfer struct {
	ID              uuid.UUID `json:"id"`
	ContractAddress string    `json:"contract_address"`
	TokenID         string    `json:"token_id"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Amount          string    `json:"amount"`
	Type            string    `json:"type"`
	TxHash          string    `json:"tx_hash"`
	LogIndex        uint64    `json:"log_index"`
}

type NFT struct {
	ID              uuid.UUID `json:"id"`
	ContractAddress string    `json:"contract_address"`
	TokenID         string    `json:"token_id"`
	OwnerAddress    string    `json:"owner_address"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	TokenURI        string    `json:"token_uri"`
	OptimizedURL    string    `json:"optimized_url"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Attributes      string    `json:"attributes"`
}

type NFTCollection struct {
	ID              uuid.UUID `json:"id"`
	ContractAddress string    `json:"contract_address"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	NumberOfNFTs    string    `json:"num_nfts"`
}
