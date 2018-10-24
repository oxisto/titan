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
	"strconv"
	"time"
)

type Skill struct {
	expireDate  *time.Time
	SkillID     int32 `json:"skillID"`
	SkillPoints int64 `json:"skillpoints"`
	Level       int32 `json:"level"`
}

type Character struct {
	expireDate      *time.Time
	CharacterID     int32            `json:"characterID"`
	CharacterName   string           `json:"name"`
	CorporationID   int32            `json:"corporationID"`
	CorporationName string           `json:"corporationName"`
	AllianceID      int32            `json:"allianceID"`
	AllianceName    string           `json:"allianceName"`
	Skills          map[string]Skill `json:"skills"` // index type needs to be string because of the flat map library
}

func (c *Character) ID() int32 {
	return c.CharacterID
}

func (c *Character) SkillLevel(skillID int32) int {
	return int(c.Skills[strconv.Itoa(int(skillID))].Level)
}

func (c *Character) ExpiresOn() *time.Time {
	return c.expireDate
}

func (c *Character) SetExpire(t *time.Time) {
	c.expireDate = t
}

func (c *Character) HashKey() string {
	return fmt.Sprintf("character:%d", c.ID())
}
