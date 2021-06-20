package model

import (
	"fmt"
	"time"
)

type Wallets struct {
	expireDate *time.Time

	CorporationID int32             `json:"corporationID"`
	Divisions     map[string]Wallet `json:"divisions"`
}

func (w *Wallets) ID() int32 {
	return w.CorporationID
}

func (w *Wallets) ExpiresOn() *time.Time {
	return w.expireDate
}

func (w *Wallets) SetExpire(t *time.Time) {
	w.expireDate = t
}

func (w *Wallets) HashKey() string {
	return fmt.Sprintf("wallets:%d", w.ID())
}

type Wallet struct {
	Division int32   `json:"division"`
	Balance  float64 `json:"balance"`
}

type JournalEntry struct {
	Amount        float64
	Balance       float64
	Date          time.Time
	Description   string
	FirstPartyID  int32  `json:"firstPartyID" db:"firstPartyID"`
	ID            int64  `json:"id" db:"id"`
	RefType       string `json:"refType" db:"refType"`
	SecondPartyID int32  `json:"secondPartyID" db:"secondPartyID"`
	CorporationID int64  `json:"corporationID" db:"corporationID"`
	Division      int32  `json:"division" db:"division"`
}

type Transaction struct {
	TransactionID int64     `json:"transactionID" db:"transactionID"`
	ClientID      int32     `json:"clientID" db:"clientID"`
	Date          time.Time `json:"date" db:"date"`
	IsBuy         bool      `json:"isBuy" db:"isBuy"`
	JournalRefID  int64     `json:"journalRefID" db:"journalRefID"`
	LocationID    int64     `json:"locationID" db:"locationID"`
	Quantity      int
	TypeID        TypeIdentifier `json:"typeID" db:"typeID"`
	UnitPrice     float64        `json:"unitPrice" db:"unitPrice"`
	CorporationID int64          `json:"corporationID" db:"corporationID"`
	Division      int32          `json:"division" db:"division"`
}
