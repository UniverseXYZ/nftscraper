package models

type Transfer struct {
	Id				uint64 `json:"id"`
	ContractAddress string `json:"contractAddress"`
	TokenId         string `json:"tokenId"`
	From            string `json:"from"`
	To              string `json:"to"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	TxHash 			string `json:"txHash"`
	LogIndex        uint64 `json:"logIndex"`
}

type Nft struct {
	// Id				uint64 `json:"id"`
	ContractAddress string `json:"contractAddress"`
	TokenId         string `json:"tokenId"`
	OwnerAddress	string `json:"ownerAddress"`
	Name			string `json:"name"`
	Symbol			string `json:"symbol"`
	OptimizedUrl	string `json:"optimizedUrl"`
	ThumbnailUrl	string `json:"thumbnailUrl"`
	Attributes		string `json:"attributes"`
}

type NftCollection struct {
	Id				uint64 `json:"id"`
	ContractAddress string `json:"contractAddress"`
	Name 			string `json:"name"`
	Symbol 			string `json:"symbol"`
	NumberOfNfts	string `json:"numberOfNfts"`			
}
