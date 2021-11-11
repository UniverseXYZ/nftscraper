
BEGIN;

CREATE TABLE IF NOT EXISTS transfer (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    contractAddress varchar(42) NOT NULL,
    tokenId varchar(100) NOT NULL,
    "from" varchar(42) NOT NULL,
    "to" varchar(42) NOT NULL,
    amount varchar(20) NOT NULL,
    "type" varchar(10) NOT NULL,
    txHash varchar(100),
    logIndex bigint
);

CREATE TABLE IF NOT EXISTS nft (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    contractAddress varchar(42) NOT NULL,
    tokenId varchar(100) NOT NULL,
    ownerAddress varchar(42) NOT NULL,
    "name" varchar(255) NOT NULL,
    symbol varchar(255) NOT NULL,
    tokenURI text NOT NULL,
    optimizedUrl text NOT NULL,
    thumbnailUrl text NOT NULL,
    attributes json
);

CREATE TABLE IF NOT EXISTS nftCollection (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    contractAddress varchar(42) NOT NULL,
    "name" varchar(255) NOT NULL,
    symbol varchar(255) NOT NULL,
    numberOfNfts varchar(20) NOT NULL
);

COMMIT;
