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

package manufacturing

import (
	"math"
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

const (
	ActivityManufacturing = 1
	ActivityInvention     = 8
)

var FacilityJobDurationBonuses = map[string]float64{
	"Station":             0,
	"Engineering Complex": -0.15,
}

var FacilityMaterialBonuses = map[string]float64{
	"Station":             0,
	"Engineering Complex": -0.01,
}

var FacilityISKBonus = map[string]float64{
	"Station":             0,
	"Engineering Complex": -0.03,
}

func CalculateModifier(modifiers map[string]float64) float64 {
	f := 1.0

	for _, mod := range modifiers {
		f += 1 * mod
	}

	return f
}

func NewManufacturing(builder *model.Character, productTypeID int32, ME int64, TE int64, facilityTax float64, object model.CachedObject) (err error) {
	manufacturing, ok := object.(*model.Manufacturing)
	if !ok {
		return errors.New("passing invalid type to NewManufacturing function")
	}

	manufacturing.ProductTypeID = productTypeID
	if manufacturing.Product, err = db.GetType(productTypeID); err != nil {
		return err
	}

	blueprint := db.GetBlueprint(productTypeID, ActivityManufacturing).Blueprint

	if blueprint.TypeID == 0 {
		return errors.New("Item cannot be manufactured")
	}

	if manufacturing.BlueprintType, err = db.GetType(blueprint.TypeID); err != nil {
		return err
	}

	industrySkillLevel := 5
	advancedIndustrySkillLevel := 5

	if builder != nil {
		industrySkillLevel = builder.SkillLevel(SkillIdIndustry)
		advancedIndustrySkillLevel = builder.SkillLevel(SkillIdAdvancedIndustry)
	}

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
		manufacturing.Invention = &model.Invention{}
		manufacturing.IsTech2 = false
		manufacturing.Runs = 1
		manufacturing.ME = ME
		manufacturing.TE = TE
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

	var materials []db.IndustryActivityMaterialResult
	if materials, err = db.GetActivityMaterials(ActivityManufacturing, blueprint, manufacturing.Runs, manufacturing.MaterialModifier); err != nil {
		return err
	}

	var skills []db.IndustryActivitySkillResult
	if skills, err = db.GetActivitySkills(ActivityManufacturing, blueprint); err != nil {
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

	eiv := 0.0

	manufacturing.Materials = map[string]model.ManufacturingMaterial{}
	for _, material := range materials {
		material.PricePerUnit = prices[material.TypeID].Sell.Percentile
		material.Cost = float64(material.Quantity) * material.PricePerUnit

		manufacturing.Materials[strconv.Itoa(int(material.TypeID))] = material.ManufacturingMaterial
		manufacturing.Costs.TotalMaterials += material.Cost

		// get market prices for Estimated Item Value (EIV) calculation
		if marketPrice, err := cache.GetMarketPrice(material.TypeID); err != nil {
			return err
		} else {
			// EIV calculation is always based on ME 0
			eiv += float64(material.RawQuantity) * marketPrice.AdjustedPrice
		}
	}

	eiv = /*math.Round(*/ eiv * float64(manufacturing.Runs) /*)*/

	manufacturing.RequiredSkills = map[string]model.ManufacturingSkill{}
	manufacturing.HasRequiredSkills = true
	for _, skill := range skills {
		skill.SkillLevel = 5
		if builder != nil {
			skill.SkillLevel = builder.SkillLevel(skill.TypeID)
		}
		skill.HasLearned = skill.SkillLevel >= skill.RequiredLevel

		if !skill.HasLearned {
			manufacturing.HasRequiredSkills = false
		}

		manufacturing.RequiredSkills[strconv.Itoa(int(skill.TypeID))] = skill.ManufacturingSkill
	}

	manufacturing.Costs.TotalJobCost = CalculateJobCost(eiv, 0.0386, FacilityISKBonus[manufacturing.Facility], facilityTax)

	manufacturing.Costs.Total = manufacturing.Costs.TotalMaterials + manufacturing.Costs.TotalJobCost
	manufacturing.Costs.PerItem = manufacturing.Costs.Total / float64(manufacturing.Product.PortionSize) / float64(manufacturing.Runs)

	if manufacturing.IsTech2 {
		manufacturing.Costs.Total += manufacturing.Invention.CostsForManufacturing
	}

	var activity db.IndustryActivityResult
	if activity, err = db.GetIndustryActivity(blueprint.TypeID, ActivityManufacturing); err != nil {
		return err
	}

	manufacturing.Time = int(math.Ceil(float64(activity.Time*manufacturing.Runs) * manufacturing.TimeModifier))
	manufacturing.SlotsUsed = manufacturing.MaxSlots
	manufacturing.ItemsPerDay = 3600.0 / float64(manufacturing.Time/manufacturing.Runs) * 24.0 * float64(manufacturing.Product.PortionSize) * float64(manufacturing.SlotsUsed)

	// revenue
	manufacturing.Revenue.Total = model.ProfitValue{
		BasedOnBuyPrice:  prices[productTypeID].Buy.Percentile * float64(manufacturing.Product.PortionSize*manufacturing.Runs),
		BasedOnSellPrice: prices[productTypeID].Sell.Percentile * float64(manufacturing.Product.PortionSize*manufacturing.Runs),
	}
	manufacturing.Revenue.PerItem = model.ProfitValue{
		BasedOnBuyPrice:  prices[productTypeID].Buy.Percentile,
		BasedOnSellPrice: prices[productTypeID].Sell.Percentile,
	}

	// profit
	manufacturing.Profit.Total = model.ProfitValue{
		BasedOnBuyPrice:  manufacturing.Revenue.Total.BasedOnBuyPrice - manufacturing.Costs.Total,
		BasedOnSellPrice: manufacturing.Revenue.Total.BasedOnSellPrice - manufacturing.Costs.Total,
	}
	manufacturing.Profit.PerItem = model.ProfitValue{
		BasedOnBuyPrice:  manufacturing.Profit.Total.BasedOnBuyPrice / float64(manufacturing.Product.PortionSize) / float64(manufacturing.Runs),
		BasedOnSellPrice: manufacturing.Profit.Total.BasedOnSellPrice / float64(manufacturing.Product.PortionSize) / float64(manufacturing.Runs),
	}
	manufacturing.Profit.PerDay = model.ProfitValue{
		BasedOnBuyPrice:  manufacturing.Profit.PerItem.BasedOnBuyPrice * manufacturing.ItemsPerDay,
		BasedOnSellPrice: manufacturing.Profit.PerItem.BasedOnSellPrice * manufacturing.ItemsPerDay,
	}

	// margin
	manufacturing.Profit.Margin = model.ProfitValue{
		BasedOnBuyPrice:  prices[productTypeID].Buy.Percentile/manufacturing.Costs.PerItem - 1,
		BasedOnSellPrice: prices[productTypeID].Sell.Percentile/manufacturing.Costs.PerItem - 1,
	}

	// other stats
	manufacturing.BuyOrderVolume = prices[productTypeID].Buy.Volume
	manufacturing.DailyBuyFactor = float64(manufacturing.BuyOrderVolume) / manufacturing.ItemsPerDay

	return nil
}

func CalculateJobCost(eiv float64, systemCostIndex float64, iskBonus float64, facilityTax float64) (jobCost float64) {
	// job cost according to the system cost index
	jobCost = eiv * systemCostIndex

	// substract ISK bonus
	jobCost += (jobCost * iskBonus)

	// apply facility tax
	jobCost += jobCost * facilityTax

	return
}
