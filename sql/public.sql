CREATE TABLE profit (
    "typeID" integer NOT NULL,
    "basedOnSellPrice" double precision,
    "basedOnBuyPrice" double precision,
    CONSTRAINT profit_pkey PRIMARY KEY ("typeID")
);

CREATE TABLE journal (
    id bigint NOT NULL,
    amount double precision,
    balance double precision,
    date timestamp(4
) WITH time zone,
    description text COLLATE pg_catalog. "default",
    "firstPartyID" integer,
    "refType" text COLLATE pg_catalog. "default",
    "secondPartyID" integer,
    CONSTRAINT journal_pkey PRIMARY KEY (id)
);

