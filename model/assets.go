/*
Copyright 2020 Christian Banse

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
