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
    date timestamp WITH time zone,
    description text COLLATE pg_catalog. "default",
    "firstPartyID" integer,
    "refType" text COLLATE pg_catalog. "default",
    "secondPartyID" integer,
    "corporationID" bigint NOT NULL,
    "division" integer NOT NULL,
    CONSTRAINT journal_pkey PRIMARY KEY (
        id
    )
);

CREATE TABLE public.transactions (
    "transactionID" bigint NOT NULL,
    "clientID" integer,
    date timestamp WITH time zone,
    "isBuy" boolean,
    "journalRefID" bigint,
    "locationID" integer,
    quantity integer,
    "typeID" integer,
    "unitPrice" double precision,
    "corporationID" bigint NOT NULL,
    "division" integer NOT NULL,
    CONSTRAINT transactions_pkey PRIMARY KEY (
        "transactionID"
    )
)
