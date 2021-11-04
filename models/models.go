package models

type Transfer struct {
	ContractAddress string `json:"contractAddress"`
	TokenId         string `json:"tokenId"` //should it be string?
	From            string `json:"from"`
	To              string `json:"to"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	TxHash 			string `json:"txHash"`
	LogIndex        uint64 `json:"logIndex"`
}

type Nft struct {
	ContractAddress string `json:"contractAddress"`
	TokenId         string `json:"tokenId"` //should it be string?
	OwnerAddress	string `json:"ownerAddress"`
	Name			string `json:"name"`
	Symbol			string `json:"symbol"`
	OptimizedUrl	string `json:"optimizedUrl"`
	ThumbnailUrl	string `json:"thumbnailUrl"`
	Attributes		string `json:"attributes"`
}

type NftCollection struct {
	Name 			string `json:"name"`
	Symbol 			string `json:"symbol"`
	Amount			string `json:"amount"` 			
}