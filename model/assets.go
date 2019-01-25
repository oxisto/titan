package model

import (
	"fmt"
	"time"
)

type CorporationAssets struct {
	expireDate *time.Time

	CorporationID int32
	Assets        map[string]Asset
}

type Asset struct {
	IsSingleton  bool
	ItemID       int64
	LocationFlag string
	LocationID   int64
	LocationType string
	Quantity     int
	TypeID       int32
}

func (c *CorporationAssets) ID() int32 {
	return c.CorporationID
}

func (c *CorporationAssets) ExpiresOn() *time.Time {
	return c.expireDate
}

func (c *CorporationAssets) SetExpire(t *time.Time) {
	c.expireDate = t
}

func (c *CorporationAssets) HashKey() string {
	return fmt.Sprintf("corporation-assets:%d", c.ID())
}
