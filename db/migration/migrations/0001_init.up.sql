BEGIN;

CREATE TABLE IF NOT EXISTS transfer (
    id uuid PRIMARY KEY NOT NULL,
    contract_addr varchar(42) NOT NULL,
    token_id varchar(100) NOT NULL,
    "from" varchar(42) NOT NULL,
    "to" varchar(42) NOT NULL,
    amount varchar(20) NOT NULL,
    "type" varchar(10) NOT NULL,
    tx_hash varchar(100),
    log_index bigint
);

CREATE TABLE IF NOT EXISTS nft (
    id uuid PRIMARY KEY NOT NULL,
    contract_addr varchar(42) NOT NULL,
    token_id varchar(100) NOT NULL,
    owner_addr varchar(42) NOT NULL,
    "name" varchar(255) NOT NULL,
    symbol varchar(255) NOT NULL,
    token_uri text NOT NULL,
    optimized_url text NOT NULL,
    thumbnail_url text NOT NULL,
    attributes jsonb
);

CREATE TABLE IF NOT EXISTS nft_collection (
    id uuid PRIMARY KEY NOT NULL,
    contract_addr varchar(42) NOT NULL,
    "name" varchar(255) NOT NULL,
    symbol varchar(255) NOT NULL,
    num_nfts varchar(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS scraper_cursor (
    "name" varchar(255) PRIMARY KEY NOT NULL,
    "last_block_num" integer NOT NULL DEFAULT 0,
    "last_tx_num" integer NOT NULL DEFAULT 0,
    "last_log_index" integer NOT NULL DEFAULT 0
);

COMMIT;
