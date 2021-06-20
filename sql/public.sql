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
);

CREATE TABLE public."industryJobs" (
    "jobID" integer NOT NULL,
    "activityID" integer NOT NULL,
	"completedCharacterID" integer NOT NULL,
	"completedDate" timestamp WITH time zone,
	"cost" double precision NOT NULL,
	"duration" integer NOT NULL,
	"endDate" timestamp WITH time zone NOT NULL,
	"facilityID" bigint NOT NULL,
	"installerID" integer NOT NULL,
	"locationID" bigint NOT NULL,
	"blueprintID" bigint NOT NULL,
	"blueprintTypeID" integer NOT NULL,
	"startDate" timestamp WITH time zone NOT NULL,
	"pauseDate" timestamp WITH time zone,
	"licensedRuns" integer NOT NULL,
	"outputLocationID" bigint NOT NULL,
	"probability" real NOT NULL,
	"productTypeID" integer NOT NULL,
	"runs" integer NOT NULL,
	"succesfulRuns" integer NOT NULL,
	"status" text COLLATE pg_catalog. "default" NOT NULL,
    CONSTRAINT industrJobs_pkey PRIMARY KEY (
        "jobID"
    )
)