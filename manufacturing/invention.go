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
	"strconv"
	"strings"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "manufacturing")
}

type SkillHolder interface {
	SkillLevel(TypeID int32) int
}

func NewInvention(blueprintTypeID int32, inventor SkillHolder) (invention *model.Invention, err error) {
	blueprint := db.GetBlueprint(blueprintTypeID, ActivityInvention).Blueprint

	invention = new(model.Invention)

	if invention.BlueprintType, err = db.GetType(blueprint.TypeID); err != nil {
		return nil, err
	}

	materials, err := db.GetActivityMaterials(ActivityInvention, blueprint, 1, 1)

	var skills []db.IndustryActivitySkillResult
	skills, err = db.GetActivitySkills(ActivityInvention, blueprint)

	typeIDs := []int32{}
	for _, material := range materials {
		typeIDs = append(typeIDs, material.TypeID)
	}

	// make sure, prices are available
	var prices map[int32]model.Price

	if prices, err = cache.GetPrices(model.JitaRegionID, typeIDs); err != nil {
		return nil, err
	}

	skillMods := []float64{}
	invention.RequiredSkills = map[string]model.ManufacturingSkill{}
	for _, skill := range skills {
		skill.SkillLevel = 5

		if inventor != nil {
			skill.SkillLevel = inventor.SkillLevel(skill.TypeID)
		}

		skill.HasLearned = skill.SkillLevel >= skill.RequiredLevel

		if strings.Contains(skill.TypeName, "Encryption") {
			skillMods = append(skillMods, float64(skill.SkillLevel)*0.0250)
		} else {
			skillMods = append(skillMods, float64(skill.SkillLevel)*0.0333)
		}

		invention.RequiredSkills[strconv.Itoa(int(skill.TypeID))] = skill.ManufacturingSkill
	}

	invention.CostsPerRun = 0
	invention.Materials = map[string]model.ManufacturingMaterial{}
	for _, material := range materials {
		material.PricePerUnit = prices[material.TypeID].Sell.Percentile
		material.Cost = float64(material.Quantity) * material.PricePerUnit

		invention.Materials[strconv.Itoa(int(material.TypeID))] = material.ManufacturingMaterial
		invention.CostsPerRun += material.Cost
	}

	// we just the the first product. the probability should be the same for all anyway
	invention.SuccessProbabilityModifiers = map[string]float64{}
	// TODO: we need a new call for this
	//invention.SuccessProbabilityModifiers["Blueprint Base Probability"] = float64(blueprint.Activities.Invention.Products[0].Probability)

	invention.SuccessProbabilityModifiers["Skills"] = 0
	for _, skillMod := range skillMods {
		invention.SuccessProbabilityModifiers["Skills"] += skillMod
	}

	invention.InventionChance = invention.SuccessProbabilityModifiers["Blueprint Base Probability"] * (1 + invention.SuccessProbabilityModifiers["Skills"])
	invention.TriesForManufacturing = 1 / invention.InventionChance
	invention.CostsForManufacturing = invention.CostsPerRun * invention.TriesForManufacturing

	return invention, nil
}
