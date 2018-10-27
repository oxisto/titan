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

package manufacturing

import (
	"math"

	"fmt"
	"time"

	"strconv"

	"errors"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
)

const (
	SkillIdIndustry         = 3380
	SkillIdAdvancedIndustry = 3388
)

var FacilityJobDurationBonuses = map[string]float64{
	"POS":                 -0.25,
	"Station":             0,
	"Engineering Complex": -0.15,
}

var FacilityMaterialBonuses = map[string]float64{
	"POS":                 -0.02,
	"Station":             0,
	"Engineering Complex": -0.01,
}

type ManufacturingMaterial struct {
	TypeID       int32               `json:"typeID" bson:"typeID"`
	Quantity     int                 `json:"quantity"`
	Name         model.LocalizedName `json:"name" bson:"name,omitempty"`
	PricePerUnit float64             `json:"pricePerUnit" bson:"pricePerUnit"`
	Cost         float64             `json:"cost"`
}

type ManufacturingSkill struct {
	TypeID        int32               `json:"typeID" bson:"typeID"`
	Name          model.LocalizedName `json:"name" bson:"name"`
	RequiredLevel int                 `json:"requiredLevel" bson:"requiredLevel"`
	SkillLevel    int                 `json:"skillLevel" bson:"skillLevel"`
	HasLearned    bool                `json:"hasLearned" bson:"hasLearned"`
}

type Manufacturing struct {
	BlueprintType                model.Type                       `json:"blueprintType" bson:"blueprintType"`
	Product                      model.Type                       `json:"product"`
	IsTech2                      bool                             `json:"isTech2" bson:"isTech2"`
	Runs                         int                              `json:"runs"`
	MaxSlots                     int                              `json:"maxSlots" bson:"maxSlots"`
	SlotsUsed                    int                              `json:"slotsUsed" bson:"slotsUsed"`
	JobDurationModifiers         map[string]float64               `json:"jobDurationModifiers" bson:"jobDurationModifiers"`
	MaterialConsumptionModifiers map[string]float64               `json:"materialConsumptionModifiers" bson:"materialConsumptionModifiers"`
	ME                           int                              `json:"me"`
	TE                           int                              `json:"te"`
	TimeModifier                 float64                          `json:"timeModifier" bson:"timeModifier"`
	MaterialModifier             float64                          `json:"materialModifier" bson:"materialModifier"`
	Materials                    map[string]ManufacturingMaterial `json:"materials"`
	RequiredSkills               map[string]ManufacturingSkill    `json:"requiredSkills" bson:"requiredSkills"`
	HasRequiredSkills            bool                             `json:"hasRequiredSkills" bson:"hasRequiredSkills"`
	Facility                     string                           `json:"facility"`
	Costs                        struct {
		TotalMaterials float64 `json:"totalMaterials" bson:"totalMaterials"`
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

func CalculateModifier(modifiers map[string]float64) float64 {
	f := 1.0

	for _, mod := range modifiers {
		f += 1 * mod
	}

	return f
}

func NewManufacturing(builder model.Character, productTypeId int32, object model.CachedObject) (err error) {
	manufacturing, ok := object.(*Manufacturing)
	if !ok {
		return errors.New("passing invalid type to NewManufacturing function")
	}

	manufacturing.Product = db.GetType(productTypeId)
	blueprint := db.GetBlueprint(productTypeId, "activities.manufacturing.products.typeID")

	if blueprint.BlueprintTypeID == 0 {
		return errors.New("Item cannot be manufactured")
	}

	manufacturing.BlueprintType = db.GetType(blueprint.BlueprintTypeID)

	industrySkillLevel := builder.SkillLevel(SkillIdIndustry)
	advancedIndustrySkillLevel := builder.SkillLevel(SkillIdAdvancedIndustry)

	if manufacturing.Product.MetaGroupID == 2 {
		manufacturing.IsTech2 = true
		if manufacturing.Invention, err = NewInvention(blueprint, builder); err != nil {
			return err
		}
		manufacturing.Runs = blueprint.MaxProductionLimit

		// TODO: decryptors
		manufacturing.ME = 2
		manufacturing.TE = 4
	} else {
		// to avoid NPE
		manufacturing.Invention = &Invention{}
		manufacturing.IsTech2 = false
		manufacturing.Runs = 1
		manufacturing.ME = 0 // TODO: from options
		manufacturing.TE = 0 // TODO: from options
	}

	manufacturing.Facility = "Engineering Complex" // TODO: from options
	manufacturing.MaxSlots = 1 + industrySkillLevel + advancedIndustrySkillLevel

	manufacturing.JobDurationModifiers = map[string]float64{}
	manufacturing.JobDurationModifiers["Skills"] = -0.04*float64(industrySkillLevel) - 0.03*float64(advancedIndustrySkillLevel-1)
	manufacturing.JobDurationModifiers["Blueprint Time Efficiency"] = -float64(manufacturing.TE) / 100
	manufacturing.JobDurationModifiers[manufacturing.Facility+" Bonus"] = FacilityJobDurationBonuses[manufacturing.Facility]

	manufacturing.MaterialConsumptionModifiers = map[string]float64{}
	manufacturing.MaterialConsumptionModifiers["Blueprint Material Efficiency"] = -float64(manufacturing.ME) / 100
	manufacturing.MaterialConsumptionModifiers[manufacturing.Facility+" Bonus"] = FacilityMaterialBonuses[manufacturing.Facility]

	manufacturing.TimeModifier = CalculateModifier(manufacturing.JobDurationModifiers)
	manufacturing.MaterialModifier = CalculateModifier(manufacturing.MaterialConsumptionModifiers)

	materials := []ManufacturingMaterial{}
	if err = db.GetActivityMaterials("manufacturing", blueprint, manufacturing.Runs, manufacturing.MaterialModifier, &materials); err != nil {
		return err
	}

	skills := []ManufacturingSkill{}
	if err = db.GetActivitySkills("manufacturing", blueprint, &skills); err != nil {
		return err
	}

	typeIDs := []int32{}
	for _, material := range materials {
		typeIDs = append(typeIDs, material.TypeID)
	}
	//typeIDs = append(typeIDs, blueprint.BlueprintTypeID)
	typeIDs = append(typeIDs, manufacturing.Product.TypeID)

	// make sure, prices are available
	var prices map[int32]model.Price

	if prices, err = cache.GetPrices(model.JitaRegionID, typeIDs); err != nil {
		return err
	}

	manufacturing.Materials = map[string]ManufacturingMaterial{}
	for _, material := range materials {
		material.PricePerUnit = prices[material.TypeID].Sell.Percentile
		material.Cost = float64(material.Quantity) * material.PricePerUnit

		manufacturing.Materials[strconv.Itoa(int(material.TypeID))] = material
		manufacturing.Costs.TotalMaterials += material.Cost
	}

	manufacturing.RequiredSkills = map[string]ManufacturingSkill{}
	manufacturing.HasRequiredSkills = true
	for _, skill := range skills {
		skill.SkillLevel = builder.SkillLevel(skill.TypeID)
		skill.HasLearned = skill.SkillLevel >= skill.RequiredLevel

		if !skill.HasLearned {
			manufacturing.HasRequiredSkills = false
		}

		manufacturing.RequiredSkills[strconv.Itoa(int(skill.TypeID))] = skill
	}

	manufacturing.Costs.Total = manufacturing.Costs.TotalMaterials // TODO: sales tax, factory tax, etc.
	manufacturing.Costs.PerItem = manufacturing.Costs.Total / float64(manufacturing.Product.PortionSize) / float64(manufacturing.Runs)

	if manufacturing.IsTech2 {
		manufacturing.Costs.Total += manufacturing.Invention.CostsForManufacturing
	}

	manufacturing.Time = int(math.Ceil(float64(blueprint.Activities.Manufacturing.Time*manufacturing.Runs) * manufacturing.TimeModifier))
	manufacturing.ItemsPerDay = 3600.0 / float64(manufacturing.Time/manufacturing.Runs) * 24.0 * float64(manufacturing.Product.PortionSize)

	if manufacturing.IsTech2 {
		manufacturing.SlotsUsed = manufacturing.MaxSlots
		manufacturing.ItemsPerDay = math.Min(manufacturing.ItemsPerDay, float64(manufacturing.Runs*manufacturing.Product.PortionSize)) * float64(manufacturing.SlotsUsed)
	} else {
		manufacturing.SlotsUsed = 1
	}

	// revenue
	manufacturing.Revenue.Total = ProfitValue{
		BasedOnBuyPrice:  prices[productTypeId].Buy.Percentile * float64(manufacturing.Product.PortionSize*manufacturing.Runs),
		BasedOnSellPrice: prices[productTypeId].Sell.Percentile * float64(manufacturing.Product.PortionSize*manufacturing.Runs),
	}
	manufacturing.Revenue.PerItem = ProfitValue{
		BasedOnBuyPrice:  prices[productTypeId].Buy.Percentile,
		BasedOnSellPrice: prices[productTypeId].Sell.Percentile,
	}

	// profit
	manufacturing.Profit.Total = ProfitValue{
		BasedOnBuyPrice:  manufacturing.Revenue.Total.BasedOnBuyPrice - manufacturing.Costs.Total,
		BasedOnSellPrice: manufacturing.Revenue.Total.BasedOnSellPrice - manufacturing.Costs.Total,
	}
	manufacturing.Profit.PerItem = ProfitValue{
		BasedOnBuyPrice:  manufacturing.Profit.Total.BasedOnBuyPrice / float64(manufacturing.Product.PortionSize) / float64(manufacturing.Runs),
		BasedOnSellPrice: manufacturing.Profit.Total.BasedOnSellPrice / float64(manufacturing.Product.PortionSize) / float64(manufacturing.Runs),
	}
	manufacturing.Profit.PerDay = ProfitValue{
		BasedOnBuyPrice:  manufacturing.Profit.PerItem.BasedOnBuyPrice * manufacturing.ItemsPerDay,
		BasedOnSellPrice: manufacturing.Profit.PerItem.BasedOnSellPrice * manufacturing.ItemsPerDay,
	}

	// margin
	manufacturing.Profit.Margin = ProfitValue{
		BasedOnBuyPrice:  prices[productTypeId].Buy.Percentile/manufacturing.Costs.PerItem - 1,
		BasedOnSellPrice: prices[productTypeId].Sell.Percentile/manufacturing.Costs.PerItem - 1,
	}

	// other stats
	manufacturing.BuyOrderVolume = prices[productTypeId].Buy.Volume
	manufacturing.DailyBuyFactor = float64(manufacturing.BuyOrderVolume) / manufacturing.ItemsPerDay

	return nil
}
