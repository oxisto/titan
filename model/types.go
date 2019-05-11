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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type TypeIdentifier int32

type Blueprint struct {
	/*ObjectID   int32 `json:"typeID" db:"_id"`
	Activities struct {
		Copying          Activity `json:"activities"`
		Invention        Activity `json:"invention"`
		Manufacturing    Activity `json:"manufacturing"`
		ResearchMaterial Activity `json:"researchMaterial"`
		ResearchTime     Activity `json:"researchTime"`
	} `json:"activities"`*/
	TypeID             int32 `json:"typeID" db:"typeID" yaml:"blueprintTypeID"`
	MaxProductionLimit int   `json:"maxProductionLimit" db:"maxProductionLimit" yaml:"maxProductionLimit"`
}

type Activity struct {
	Time      int             `json:"time"`
	Materials []Material      "materials,omitempty"
	Products  []Product       "products,omitempty"
	Skills    []RequiredSkill "skills,omitempty"
}

type Material struct {
	TypeID   int32 `json:"typeID" db:"typeID" yaml:"typeID"`
	Quantity int   `json:"quantity"`
}

type Product struct {
	Probability float32
	Quantity    int
	TypeID      int32 `json:"typeID" db:"typeID" yaml:"typeID"`
}

type RequiredSkill struct {
	TypeID int32 `json:"typeID" db:"typeID" yaml:"typeID"`
	Level  int
}

type Type struct {
	TypeID        int32    `json:"typeID" db:"typeID"`
	BasePrice     *float32 `json:"basePrice" db:"basePrice" yaml:"basePrice"`
	Description   string   `json:"description,omitempty" db:"description,omitempty"`
	GroupID       int32    `json:"groupID" db:"groupID" yaml:"groupID"`
	MarketGroupID *int32   `json:"marketGroupID" db:"marketGroupID" yaml:"marketGroupID"`
	Group                  // not part of SDE, but joined from postgres
	TypeName      string   `json:"typeName" db:"typeName"`
	Mass          float64  `db:"mass"`
	Capacity      float64
	IconID        *int32  `json:"iconID" db:"iconID"`
	SoundID       *int32  `json:"soundID" db:"soundID"`
	GraphicID     int32   `json:"graphicID" db:"graphicID"`
	MetaGroupID   int32   `json:"metaGroupID" db:"metaGroupID" yaml:"metaGroupID"`
	PortionSize   int     `json:"portionSize" db:"portionSize" yaml:"portionSize"`
	Published     bool    `json:"published"`
	RaceID        *int32  `json:"raceID" db:"raceID" yaml:"raceID"`
	Volume        float64 `json:"volume"`
}

func (t Type) ID() int32 {
	return t.TypeID
}

func (t Type) ExpiresOn() *time.Time {
	return nil
}

func (t Type) SetExpire(time *time.Time) {

}

func (t Type) HashKey() string {
	return fmt.Sprintf("type:%d", t.ID())
}

type Group struct {
	GroupID    int32  `json:"groupID" db:"groupID"`
	CategoryID int32  `json:"categoryID" db:"categoryID" yaml:"categoryID"`
	GroupName  string `json:"groupName" db:"groupName"`
	Published  bool   `json:"published"`
}

type Category struct {
	CategoryID   int32  `json:"categoryID" db:"categoryID"`
	CategoryName string `json:"categoryName" db:"categoryName"`
	Published    bool   `json:"published"`
	IconID       *int32 `json:"iconID" db:"iconID"`
}

type LocalizedName struct {
	DE string `json:"de" db:"de"`
	EN string `json:"en" db:"en"`
	FR string `json:"fr" db:"fr"`
	JA string `json:"ja" db:"ja"`
	RU string `json:"ru" db:"ru"`
	ZH string `json:"zh" db:"zh"`
}

func (p LocalizedName) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *LocalizedName) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var name LocalizedName
	err := json.Unmarshal(source, &name)
	if err != nil {
		return err
	}

	*p = name

	return nil
}

type MetaGroup struct {
	MetaGroupID  int32 "metaGroupID"
	ParentTypeID int32 "parentTypeID"
	TypeID       int32 "typeID"
}
