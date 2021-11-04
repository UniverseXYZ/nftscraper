
CREATE TABLE IF NOT EXISTS transfer (
    contractAddress varchar(42) NOT NULL,
    tokenId bigint NOT NULL,
    `from` varchar(42) NOT NULL,
    `to` varchar(42) NOT NULL,
    amount bigint NOT NULL,
    type varchar(10) NOT NULL,
    txHash varchar(66),
    logIndex bigint
);