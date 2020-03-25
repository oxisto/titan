package model

import (
	"fmt"
	"time"
)

type ManufacturingMaterial struct {
	// TypeID is the id of the type that is beeing manufactured
	TypeID       int32   `json:"typeID" db:"typeID"`
	Quantity     int     `json:"quantity"`
	RawQuantity  int     `json:"rawQuantity" db:"rawQuantity"`
	TypeName     string  `json:"typeName" db:"typeName"`
	PricePerUnit float64 `json:"pricePerUnit" db:"pricePerUnit"`
	Cost         float64 `json:"cost"`
}

type ManufacturingSkill struct {
	TypeID        int32  `json:"typeID" db:"typeID"`
	TypeName      string `json:"typeName" db:"typeName"`
	RequiredLevel int    `json:"requiredLevel" db:"requiredLevel"`
	SkillLevel    int    `json:"skillLevel" db:"skillLevel"`
	HasLearned    bool   `json:"hasLearned" db:"hasLearned"`
}

type IndustryActivity struct {
	TypeID     int32 `json:"typeID" db:"typeID"`
	Time       int   `json:"time"`
	ActivityID int32 `json:"activityID" db:"activityID"`
}

type Manufacturing struct {
	BlueprintType                Type                             `json:"blueprintType" bson:"blueprintType"`
	Product                      Type                             `json:"product"`
	ProductTypeID                int32                            `json:"productTypeID"`
	IsTech2                      bool                             `json:"isTech2" bson:"isTech2"`
	Runs                         int                              `json:"runs"`
	MaxSlots                     int                              `json:"maxSlots" bson:"maxSlots"`
	SlotsUsed                    int                              `json:"slotsUsed" bson:"slotsUsed"`
	JobDurationModifiers         map[string]float64               `json:"jobDurationModifiers" bson:"jobDurationModifiers"`
	MaterialConsumptionModifiers map[string]float64               `json:"materialConsumptionModifiers" bson:"materialConsumptionModifiers"`
	ME                           int64                            `json:"me"`
	TE                           int64                            `json:"te"`
	TimeModifier                 float64                          `json:"timeModifier" bson:"timeModifier"`
	MaterialModifier             float64                          `json:"materialModifier" bson:"materialModifier"`
	Materials                    map[string]ManufacturingMaterial `json:"materials"`
	RequiredSkills               map[string]ManufacturingSkill    `json:"requiredSkills" bson:"requiredSkills"`
	HasRequiredSkills            bool                             `json:"hasRequiredSkills" bson:"hasRequiredSkills"`
	Facility                     string                           `json:"facility"`
	Costs                        struct {
		TotalMaterials float64 `json:"totalMaterials" bson:"totalMaterials"`
		TotalJobCost   float64 `json:"totalJobCost" bson:"totalJobCost"`
		Total          float64 `json:"total"`
		PerItem        float64 `json:"perItem" bson:"perItem"`
	} `json:"costs"`
	Revenue struct {
		Total   ProfitValue `json:"total"`
		PerItem ProfitValue `json:"perItem" bson:"perItem"`
	} `json:"revenue"`
	Profit struct {
		Total   ProfitValue `json:"total"`
		PerItem ProfitValue `json:"perItem" bson:"perItem"`
		PerDay  ProfitValue `json:"perDay" bson:"perDay"`
		Margin  ProfitValue `json:"margin"`
	} `json:"profit"`
	BuyOrderVolume int        `json:"buyOrderVolume" bson:"buyOrderVolume"`
	DailyBuyFactor float64    `json:"dailyBuyFactor" bson:"dailyBuyFactor"`
	Time           int        `json:"time"`
	ItemsPerDay    float64    `json:"itemsPerDay" bson:"itemsPerDay"`
	Invention      *Invention `json:"invention"`
}

func (m Manufacturing) ID() int32 {
	return m.Product.TypeID
}

func (m Manufacturing) ExpiresOn() *time.Time {
	return nil
}

func (m Manufacturing) SetExpire(t *time.Time) {

}

func (m Manufacturing) HashKey() string {
	return fmt.Sprintf("manufacturing:%d", m.ID())
}

type ProfitValue struct {
	BasedOnBuyPrice  float64 `json:"basedOnBuyPrice" bson:"basedOnBuyPrice"`
	BasedOnSellPrice float64 `json:"basedOnSellPrice" bson:"basedOnSellPrice"`
}

type Invention struct {
	BlueprintType               Type                             `json:"blueprintType" bson:"blueprintType"`
	CostsPerInvention           int                              `json:"costsPerInvention" bson:"costsPerInvention"`
	DecryptorTypeID             int32                            `json:"decryptorTypeID" bson:"decryptorTypeID"`
	Materials                   map[string]ManufacturingMaterial `json:"materials"`
	RequiredSkills              map[string]ManufacturingSkill    `json:"requiredSkills" bson:"requiredSkills"`
	SuccessProbabilityModifiers map[string]float64               `json:"successProbabilityModifiers" bson:"successProbabilityModifiers"`
	CostsPerRun                 float64                          `json:"costsPerRun" bson:"costsPerRun"`
	InventionChance             float64                          `json:"inventionChance" bson:"inventionChance"`
	TriesForManufacturing       float64                          `json:"triesForManufacturing" bson:"triesForManufacturing"`
	CostsForManufacturing       float64                          `json:"costsForManufacturing" bson:"costsForManufacturing"`
}
