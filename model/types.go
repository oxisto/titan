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

type Blueprint struct {
	ObjectID   int32 `json:"typeID" bson:"_id"`
	Activities struct {
		Copying          Activity `json:"activities"`
		Invention        Activity `json:"activities"`
		Manufacturing    Activity `json:"activities"`
		ResearchMaterial Activity "research_material"
		ResearchTime     Activity "research_time"
	} `json:"activities"`
	BlueprintTypeID    int32 `json:"blueprintTypeID" bson:"blueprintTypeID" yaml:"blueprintTypeID"`
	MaxProductionLimit int   `json:"maxProductionLimit" bson:"maxProductionLimit" yaml:"maxProductionLimit"`
}

type Activity struct {
	Time      int
	Materials []Material      "materials,omitempty"
	Products  []Product       "products,omitempty"
	Skills    []RequiredSkill "skills,omitempty"
}

type Material struct {
	TypeID   int32 `json:"typeID" bson:"typeID" yaml:"typeID"`
	Quantity int   `json:"quantity"`
}

type Product struct {
	Probability float32
	Quantity    int
	TypeID      int32 `json:"typeID" bson:"typeID" yaml:"typeID"`
}

type RequiredSkill struct {
	TypeID int32 `json:"typeID" bson:"typeID" yaml:"typeID"`
	Level  int
}

type Type struct {
	TypeID      int32         `json:"typeID" bson:"_id"`
	BasePrice   int           `json:"basePrice" bson:"basePrice" yaml:"basePrice"`
	Description LocalizedName `json:"description,omitempty" bson:"description,omitempty"`
	GroupID     int32         `json:"groupID" bson:"groupID" yaml:"groupID"`
	Group       Group         `json:"group" bson:"group,omitempty"` // not part of SDE, but joined from mongo
	Name        LocalizedName `json:"name"`
	MetaGroupID int32         `json:"metaGroupID" bson:"metaGroupID" yaml:"metaGroupID"`
	PortionSize int           `json:"portionSize" bson:"portionSize" yaml:"portionSize"`
	Published   bool          `json:"published"`
	RaceID      int32         `json:"raceID" bson:"raceID" yaml:"raceID"`
	Volume      float32       `json:"volume"`
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
	GroupID    int32         `json:"groupID" bson:"_id"`
	CategoryID int32         `json:"categoryID" bson:"categoryID" yaml:"categoryID"`
	Name       LocalizedName `json:"name"`
	Published  bool          `json:"published"`
}

type Category struct {
	CategoryID int32         `json:"categoryID" bson:"_id"`
	Name       LocalizedName `json:"name"`
	Published  bool          `json:"published"`
}

type LocalizedName struct {
	DE string `json:"de" bson:"de"`
	EN string `json:"en" bson:"en"`
	FR string `json:"fr" bson:"fr"`
	JA string `json:"ja" bson:"ja"`
	RU string `json:"ru" bson:"ru"`
	ZH string `json:"zh" bson:"zh"`
}

type MetaGroup struct {
	MetaGroupID  int32 "metaGroupID"
	ParentTypeID int32 "parentTypeID"
	TypeID       int32 "typeID"
}
