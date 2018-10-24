package model

import (
	"fmt"
	"time"
)

const (
	JitaRegionID = 10000002
)

type Price struct {
	expireDate *time.Time
	TypeID     int32
	Buy        PriceData
	Sell       PriceData
}

type PriceData struct {
	WeightedAverage float64 `json:"weightedAverage,string"`
	Max             float64 `json:"max,string"`
	Min             float64 `json:"min,string"`
	StdDev          float64 `json:"stddev,string"`
	Median          float64 `json:"median,string"`
	Volume          int     `json:"volume,string"`
	OrderCount      int     `json:"orderCount,string"`
	Percentile      float64 `json:"percentile,string"`
}

func (c *Price) ID() int32 {
	return c.TypeID
}

func (c *Price) ExpiresOn() *time.Time {
	return c.expireDate
}

func (c *Price) SetExpire(t *time.Time) {
	c.expireDate = t
}

func (c *Price) HashKey() string {
	return fmt.Sprintf("price:%d", c.ID())
}
