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

package db

import (
	"fmt"
	"math"

	"github.com/oxisto/titan/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var pdb *sqlx.DB
var log *logrus.Entry
var host string

func init() {
	log = logrus.WithField("component", "cache")
}

func InitPostgreSQL(h string) {
	host = h
	pdb, _ = sqlx.Connect("postgres", fmt.Sprintf("postgres://postgres@%s/titan?sslmode=disable", host))

	log.Infof("Using PostgreSQL @ %s", h)
}

type Profit struct {
	BasedOnBuyPrice  *float64 `json:"basedOnBuyPrice" db:"basedOnBuyPrice"`
	BasedOnSellPrice *float64 `json:"basedOnSellPrice" db:"basedOnSellPrice"`
}

type ProductTypeResult struct {
	TypeID     int    `json:"typeID" db:"typeID"`
	TypeName   string `json:"typeName" db:"typeName"`
	CategoryID int    `json:"categoryID" db:"categoryID"`
	Profit
	Costs struct {
		Total float64 `json:"total"`
	} `json:"costs"`
	HasRequiredSkills bool `json:"hasRequiredSkills"`
}

type IndustryActivityResult struct {
	model.IndustryActivity
}

type IndustryActivityMaterialResult struct {
	model.ManufacturingMaterial
}

type IndustryActivitySkillResult struct {
	model.ManufacturingSkill
}

type IndustryActivityProbabilityResult struct {
	TypeID        int32   `json:"typeID" db:"typeID"`
	ActivityID    int32   `json:"activityID" db:"activityID"`
	ProductTypeID int32   `json:"productTypeID" db:"productTypeID"`
	Probability   float64 `json:"probability" db:"probability"`
}

type BlueprintResult struct {
	model.Blueprint
}

type SearchOptions struct {
	CategoryIDs           map[int]bool
	SortByField           string
	SortByDirection       string
	NameFilter            string
	MaxProductionCosts    float64
	MetaGroupID           int
	Offset                int
	Limit                 int
	HasRequiredSkillsOnly bool
}

func NewSearchOptions() *SearchOptions {
	options := &SearchOptions{}
	options.SortByField = "basedOnSellPrice"
	options.SortByDirection = "DESC"
	options.MaxProductionCosts = math.MaxInt32
	options.Limit = 100
	options.Offset = 0

	return options
}

func UpdateProfit(m model.Manufacturing) {
	log.Debugf("Updating profit for %s (%d)...", m.Product.TypeName, m.Product.TypeID)

	_, err := pdb.Exec(`INSERT INTO profit ("typeID", "basedOnSellPrice", "basedOnBuyPrice")
        VALUES ($1, $2, $3) ON CONFLICT ("typeID")
        DO
        UPDATE
        SET
            "basedOnSellPrice" = excluded. "basedOnSellPrice",
            "basedOnBuyPrice" = excluded. "basedOnBuyPrice"
`, m.Product.TypeID, m.Profit.PerDay.BasedOnSellPrice, m.Profit.PerDay.BasedOnBuyPrice)

	if err != nil {
		log.Printf("Could not update profit: %v", err)
	}
}

func GetMaterialTypeIDs(activityID model.IndustryActivityID) []int32 {
	typeIDs := []int32{}

	pdb.Select(&typeIDs, `SELECT DISTINCT
    "typeID"
FROM
    evesde. "industryActivityMaterials"
WHERE
    "activityID" = $1
`, activityID)

	return typeIDs
}

func GetIndustryActivity(typeID int32, activityID model.IndustryActivityID) (IndustryActivityResult, error) {
	activity := IndustryActivityResult{}

	// TODO: directly join materials?

	err := pdb.Get(&activity, `SELECT
    "industryActivity".*
FROM
    evesde. "industryActivity"
WHERE
    "typeID" = $1
    AND "activityID" = $2
`, typeID, activityID)

	return activity, err
}

func GetCategories() ([]model.Category, error) {
	categories := []model.Category{}

	err := pdb.Select(&categories, `SELECT
    *
FROM
    evesde. "invCategories"
WHERE
    published = TRUE
`)

	return categories, err
}

func GetActivityMaterials(activityID model.IndustryActivityID, blueprint model.Blueprint, runs int, materialModifier float64) ([]IndustryActivityMaterialResult, error) {
	materials := []IndustryActivityMaterialResult{}

	err := pdb.Select(&materials, `SELECT
    "invTypes"."typeID",
    "invTypes"."typeName",
	CAST(CEIL(quantity * $3 * $4::double precision) AS integer) AS quantity,
	"quantity" AS "rawQuantity"
FROM
    evesde. "industryActivityMaterials"
    JOIN evesde. "invTypes" ON ("industryActivityMaterials"."materialTypeID" = "invTypes"."typeID")
WHERE
    "industryActivityMaterials"."typeID" = $1
    AND "industryActivityMaterials"."activityID" = $2
ORDER BY
    "invTypes"."typeName"
`, blueprint.TypeID, activityID, runs, materialModifier)

	return materials, err
}

func GetActivityProbablitities(activityID model.IndustryActivityID, blueprintTypeID int32) ([]IndustryActivityProbabilityResult, error) {
	var result = []IndustryActivityProbabilityResult{}

	err := pdb.Select(&result, `SELECT * FROM 
	evesde."industryActivityProbabilities"
WHERE 
	"typeID" = $1
	AND	"activityID" = $2
`, blueprintTypeID, activityID)

	return result, err
}

func GetActivityProbablity(activityID model.IndustryActivityID, blueprintTypeID int32, productTypeID int32) (IndustryActivityProbabilityResult, error) {
	var result = IndustryActivityProbabilityResult{}

	err := pdb.Get(&result, `SELECT * FROM 
	evesde."industryActivityProbabilities"
WHERE 
	"typeID" = $1
	AND	"activityID" = $2
	AND "productTypeID" = $3
`, blueprintTypeID, activityID, productTypeID)

	return result, err
}

func GetActivitySkills(activityID model.IndustryActivityID, blueprint model.Blueprint) ([]IndustryActivitySkillResult, error) {
	skills := []IndustryActivitySkillResult{}

	err := pdb.Select(&skills, `SELECT
    "invTypes"."typeID",
    "invTypes"."typeName",
	"industryActivitySkills".level AS "requiredLevel"
FROM
    evesde. "industryActivitySkills"
    JOIN evesde. "invTypes" ON ("industryActivitySkills"."skillID" = "invTypes"."typeID")
WHERE
    "industryActivitySkills"."typeID" = $1
    AND "industryActivitySkills"."activityID" = $2
ORDER BY
    "invTypes"."typeName"
`, blueprint.TypeID, activityID)

	return skills, err
}

func GetBlueprint(activityID model.IndustryActivityID, productTypeID int32) BlueprintResult {
	blueprint := BlueprintResult{}

	pdb.Get(&blueprint, `SELECT
    "industryBlueprints".*
FROM
    evesde. "industryActivityProducts"
    JOIN evesde. "industryBlueprints" USING ("typeID")
WHERE
    "activityID" = $1
    AND "productTypeID" = $2
`, activityID, productTypeID)

	return blueprint
}

func GetType(typeID int32) (*model.Type, error) {
	t := model.Type{}

	err := pdb.Get(&t, `SELECT
    "invTypes".*,
    "invGroups"."categoryID",
    "invGroups"."groupName",
	"invMetaTypes"."metaGroupID"
FROM
    evesde. "invTypes"
    JOIN evesde. "invGroups" USING ("groupID")
	LEFT JOIN evesde. "invMetaTypes" USING ("typeID")
WHERE
    "typeID" = $1
`, typeID)

	return &t, err
}

func GetProductTypeIDs() ([]int32, error) {
	types := []int32{}

	err := pdb.Select(&types, `SELECT
    "invTypes"."typeID"
FROM
    evesde. "industryActivityProducts"
    JOIN evesde. "invTypes" ON ("invTypes"."typeID" = "productTypeID")
    LEFT JOIN evesde. "invMetaTypes" ON ("invMetaTypes"."typeID" = "invTypes"."typeID")
WHERE
    "activityID" = 1
    AND published = TRUE
    AND ("metaGroupID" IS NULL
        OR "metaGroupID" IN (1, 2))
`)

	return types, err
}

func GetTech1BlueprintIDs() []int32 {
	types := []int32{}

	pdb.Select(&types, `SELECT
    "industryActivityProducts"."typeID"
FROM
    evesde. "industryActivityProducts"
    JOIN evesde. "invTypes" ON ("invTypes"."typeID" = "productTypeID")
    LEFT JOIN evesde. "invMetaTypes" ON ("invMetaTypes"."typeID" = "invTypes"."typeID")
WHERE
    "activityID" = 1
    AND published = TRUE
    AND ("metaGroupID" IS NULL
        OR "metaGroupID" IN (1))`)

	return types
}

func GetProductTypes(options *SearchOptions) ([]ProductTypeResult, error) {
	types := []ProductTypeResult{}

	if options == nil {
		options = NewSearchOptions()
	}

	err := pdb.Select(&types, `SELECT
    "invTypes"."typeID",
    "invTypes"."typeName",
    "invGroups"."categoryID",
    profit. "basedOnBuyPrice",
    profit. "basedOnSellPrice"
FROM
    evesde. "industryActivityProducts"
    JOIN evesde. "invTypes" ON ("invTypes"."typeID" = "productTypeID")
    LEFT JOIN evesde. "invMetaTypes" ON ("invMetaTypes"."typeID" = "invTypes"."typeID")
    LEFT JOIN evesde. "invGroups" USING ("groupID")
    LEFT JOIN profit ON ("invTypes"."typeID" = profit. "typeID")
WHERE
    "activityID" = 1
    AND "invTypes".published = TRUE
    AND ("metaGroupID" IS NULL
        OR "metaGroupID" IN (1, 2))
    AND "typeName" ILIKE $2
ORDER BY
    "`+options.SortByField+`" DESC NULLS LAST, "typeName"
LIMIT $1`, options.Limit, "%"+options.NameFilter+"%")

	return types, err
}
