package manufacturing

import (
	"strconv"
	"strings"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
)

type Invention struct {
	BlueprintType               model.Type                       `json:"blueprintType" bson:"blueprintType"`
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

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "manufacturing")
}

func NewInvention(tech2Blueprint model.Blueprint, inventor model.Character) (invention *Invention, err error) {
	blueprint := db.GetBlueprint(tech2Blueprint.BlueprintTypeID, "activities.invention.products.typeID")

	invention = &Invention{}

	if blueprint.ObjectID == 0 {
		log.Warnf("Why?? TypeID %d", blueprint.ObjectID)
		return
	}

	invention.BlueprintType = db.GetType(blueprint.BlueprintTypeID)

	materials := []ManufacturingMaterial{}
	db.GetActivityMaterials("invention", blueprint, 1, 1, &materials)

	skills := []ManufacturingSkill{}
	db.GetActivitySkills("invention", blueprint, &skills)

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
	invention.RequiredSkills = map[string]ManufacturingSkill{}
	for _, skill := range skills {
		skill.SkillLevel = inventor.SkillLevel(skill.SkillID)
		skill.HasLearned = skill.SkillLevel >= skill.RequiredLevel

		if strings.Contains(skill.SkillName.EN, "Encryption") {
			skillMods = append(skillMods, float64(skill.SkillLevel)*0.0250)
		} else {
			skillMods = append(skillMods, float64(skill.SkillLevel)*0.0333)
		}

		invention.RequiredSkills[strconv.Itoa(int(skill.SkillID))] = skill
	}

	invention.CostsPerRun = 0
	invention.Materials = map[string]ManufacturingMaterial{}
	for _, material := range materials {
		material.PricePerUnit = prices[material.TypeID].Sell.Percentile
		material.Cost = float64(material.Quantity) * material.PricePerUnit

		invention.Materials[strconv.Itoa(int(material.TypeID))] = material
		invention.CostsPerRun += material.Cost
	}

	// we just the the first product. the probability should be the same for all anyway
	invention.SuccessProbabilityModifiers = map[string]float64{}
	invention.SuccessProbabilityModifiers["Blueprint Base Probability"] = float64(blueprint.Activities.Invention.Products[0].Probability)

	invention.SuccessProbabilityModifiers["Skills"] = 0
	for _, skillMod := range skillMods {
		invention.SuccessProbabilityModifiers["Skills"] += skillMod
	}

	invention.InventionChance = invention.SuccessProbabilityModifiers["Blueprint Base Probability"] * (1 + invention.SuccessProbabilityModifiers["Skills"])
	invention.TriesForManufacturing = 1 / invention.InventionChance
	invention.CostsForManufacturing = invention.CostsPerRun * invention.TriesForManufacturing

	return invention, nil
}
