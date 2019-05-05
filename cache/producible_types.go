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

package cache

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
)

type ProfitValue struct {
	BasedOnBuyPrice  float64 `json:"basedOnBuyPrice" bson:"basedOnBuyPrice"`
	BasedOnSellPrice float64 `json:"basedOnSellPrice" bson:"basedOnSellPrice"`
}

type ProductTypeResult struct {
	TypeID     int    `json:"typeID" db:"typeID"`
	TypeName   string `json:"name"`
	CategoryID int    `json:"categoryID"`
	Profit     struct {
		PerDay ProfitValue `json:"perDay"`
	} `json:"profit"`
	Costs struct {
		Total float64 `json:"total"`
	} `json:"costs"`
	HasRequiredSkills bool `json:"hasRequiredSkills"`
}

func GetProductTypes(options *db.SearchOptions, builder model.Character) ([]ProductTypeResult, error) {
	s := redis.Sort{
		By: "manufacturing:*->" + options.SortByField,
		Get: []string{
			"#",
			"manufacturing:*->Product.Name.EN",
			"manufacturing:*->Product.Group.CategoryID",
			"manufacturing:*->Profit.PerDay.BasedOnSellPrice",
			"manufacturing:*->Profit.PerDay.BasedOnBuyPrice",
			"manufacturing:*->Costs.Total",
			"manufacturing:*->HasRequiredSkills",
		},
		Order: options.SortByDirection,
	}

	results := []string{}
	results, err := cache.Sort("productTypeIDs", &s).Result()
	if err != nil {
		return nil, err
	}

	columns := len(s.Get)
	numTypes := len(results) / columns

	types := []ProductTypeResult{}

	for i := 0; i < numTypes; i++ {
		t := ProductTypeResult{}
		t.TypeID, err = strconv.Atoi(results[i*columns])
		t.TypeName = results[i*columns+1]
		t.CategoryID, err = strconv.Atoi(results[i*columns+2])
		t.Profit.PerDay.BasedOnSellPrice, err = strconv.ParseFloat(results[i*columns+3], 64)
		t.Profit.PerDay.BasedOnBuyPrice, err = strconv.ParseFloat(results[i*columns+4], 64)
		t.Costs.Total, err = strconv.ParseFloat(results[i*columns+5], 64)
		t.HasRequiredSkills, err = strconv.ParseBool(results[i*columns+6])

		if err != nil {
			continue
		}

		if options.NameFilter != "" && !strings.Contains(strings.ToLower(t.TypeName), strings.ToLower(options.NameFilter)) {
			continue
		}

		if len(options.CategoryIDs) > 0 && !options.CategoryIDs[t.CategoryID] {
			continue
		}

		if options.MaxProductionCosts != 0 && options.MaxProductionCosts < t.Costs.Total {
			continue
		}

		if options.HasRequiredSkillsOnly && !t.HasRequiredSkills {
			continue
		}

		types = append(types, t)
	}

	// we need to offset/limit on the types instead of Redis since we are filtering after retrieving from Redis
	var limit int
	switch {
	case options.Limit == -1:
		limit = len(types)
	case options.Limit > len(types):
		limit = len(types)
	default:
		limit = options.Limit
	}

	return types[options.Offset:limit], nil
}

/*func GetProductTypeIDs() ([]int32, error) {
	exists, err := cache.Exists("productTypeIDs").Result()
	if err != nil {
		return nil, err
	}
	if exists != 1 {
		log.Info("Fetching producible types from DB...")

		typeIDs, err := db.GetProductTypeIDs()
		if err != nil {
			return nil, err
		}

		log.Debug("Done")

		for _, t := range typeIDs {
			cache.SAdd("productTypeIDs", t)
		}

		return typeIDs, nil
	} else {
		typeIDs := []int32{}

		err = cache.SMembers("productTypeIDs").ScanSlice(&typeIDs)

		return typeIDs, err
	}
}*/
