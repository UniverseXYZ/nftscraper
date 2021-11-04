
CREATE TABLE IF NOT EXISTS `transfer` (
    id bigserial PRIMARY KEY,
    contractAddress varchar(42) NOT NULL,
    tokenId string(20) NOT NULL,
    `from` varchar(42) NOT NULL,
    `to` varchar(42) NOT NULL,
    amount varchar(20) NOT NULL,
    `type` varchar(10) NOT NULL,
    txHash varchar(100),
    logIndex bigint
);

CREATE TABLE IF NOT EXISTS `nft` (
    -- id bigserial PRIMARY KEY,
    contractAddress varchar(42) NOT NULL,
    tokenId string(20) NOT NULL,
    ownerAddress varchar(42) NOT NULL,
    `name` varchar(255) NOT NULL,
    symbol varchar(255) NOT NULL,
    optimizedUrl text NOT NULL,
    thumbnailUrl text NOT NULL,
    attributes json,
    PRIMARY KEY(contractAddress, tokenId)
);

CREATE TABLE IF NOT EXISTS `nftCollection` (
    id bigserial PRIMARY KEY,
    contractAddress varchar(42) NOT NULL,
    `name` varchar(255) NOT NULL,
    symbol varchar(255) NOT NULL,
    numberOfNfts varchar(20) NOT NULL,
    UNIQUE(contractAddress)
);