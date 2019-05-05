CREATE TABLE profit (
    "typeID" integer NOT NULL,
    "basedOnSellPrice" double precision,
    "basedOnBuyPrice" double precision,
    CONSTRAINT profit_pkey PRIMARY KEY ("typeID"))
