package model

import "time"

type JournalEntry struct {
	Amount        float64
	Balance       float64
	Date          time.Time
	Description   string
	FirstPartyID  int32  `json:"firstPartyID" db:"firstPartyID"`
	ID            int64  `json:"id" db:"id"`
	RefType       string `json:"refType" db:"refType"`
	SecondPartyID int32  `json:"secondPartyID" db:"secondPartyID"`
}

type Transaction struct {
	TransactionID int64 `json:"transactionID" db:"transactionID"`
	ClientID      int32 `json:"clientID" db:"clientID"`
	Date          time.Time
	IsBuy         bool  `json:"isBuy" db:"isBuy"`
	JournalRefID  int64 `json:"journalRefID" db:"journalRefID"`
	LocationID    int64 `json:"locationID" db:"locationID"`
	Quantity      int
	TypeID        TypeIdentifier `json:"typeID" db:"typeID"`
	UnitPrice     float64        `json:"unitPrice" db:"unitPrice"`
}
