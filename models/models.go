package models

import "github.com/google/uuid"

type Transfer struct {
	ID				uuid.UUID `json:"id"`
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenId"`
	From            string `json:"from"`
	To              string `json:"to"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	TxHash 			string `json:"txHash"`
	LogIndex        uint64 `json:"logIndex"`
}

type NFT struct {
	ID				uuid.UUID `json:"id"`
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenId"`
	OwnerAddress	string `json:"ownerAddress"`
	Name			string `json:"name"`
	Symbol			string `json:"symbol"`
	TokenURI		string `json:"tokenUri"`
	OptimizedURL	string `json:"optimizedUrl"`
	ThumbnailURL	string `json:"thumbnailUrl"`
	Attributes		string `json:"attributes"`
}

type NFTCollection struct {
	ID				uuid.UUID `json:"id"`
	ContractAddress string `json:"contractAddress"`
	Name 			string `json:"name"`
	Symbol 			string `json:"symbol"`
	NumberOfNFTs	string `json:"numberOfNfts"`			
}
