/*
Copyright 2018 Christian Banse

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
