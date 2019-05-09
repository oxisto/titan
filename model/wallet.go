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
