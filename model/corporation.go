package model

import (
	"fmt"
	"time"
)

type Corporation struct {
	expireDate      *time.Time
	CorporationID   int32            `json:"corporationID"`
	CorporationName string           `json:"Name"`
	AllianceID      int32            `json:"allianceID"`
	CEOID           int32            `json:"CEOID"`
	Ticker          string           `json:"ticker"`
	Members         map[string]int32 `json:"members"`
}

func (c *Corporation) ID() int32 {
	return c.CorporationID
}

func (c *Corporation) ExpiresOn() *time.Time {
	return c.expireDate
}

func (c *Corporation) SetExpire(t *time.Time) {
	c.expireDate = t
}

func (c *Corporation) HashKey() string {
	return fmt.Sprintf("corporation:%d", c.ID())
}
