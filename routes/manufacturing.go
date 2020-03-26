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

package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/manufacturing"
	"github.com/oxisto/titan/model"

	"github.com/oxisto/go-httputil"
)

const (
	QueryParamCategoryIDs           = "categoryIDs"
	QueryParamNameFilter            = "nameFilter"
	QueryParamSortBy                = "sortBy"
	QueryParamMaxProductionCosts    = "maxProductionCosts"
	QueryParamHasRequiredSkillsOnly = "hasRequiredSkillsOnly"
	QueryParamME                    = "ME"
	QueryParamTE                    = "TE"
	QueryParamFacilityTax           = "facilityTax"

	RouteVarsTypeID = "typeID"

	SeparatorCategoryIDs = ","
	SeparatorSortBy      = ":"
)

func GetManufacturingCategories(c *gin.Context) {
	categories, err := db.GetCategories()

	httputil.JSON(c, http.StatusOK, categories, err)
}

type ManufacturingResponse struct {
	TypeID        int32
	Type          model.Type
	Manufacturing *model.Manufacturing
}

func GetManufacturing(c *gin.Context) {
	var (
		typeID int64
		err    error
	)

	ME, err := httputil.IntQuery(c, QueryParamME)
	TE, err := httputil.IntQuery(c, QueryParamTE)
	facilityTax, err := strconv.ParseFloat(c.Query(QueryParamFacilityTax), 64)

	character := c.Value(CharacterContext).(*model.Character)

	if typeID, err = httputil.IntParam(c, "id"); err != nil {
		httputil.JSON(c, http.StatusBadRequest, nil, err)
		return
	}

	log.Debugf("Calculating manufacturing information for typeID %d...", typeID)

	resp := ManufacturingResponse{}
	m := model.Manufacturing{}

	// calculate it fresh
	if err = manufacturing.NewManufacturing(character, int32(typeID), ME, TE, facilityTax, &m); err == nil {
		resp.Manufacturing = &m
	}

	httputil.JSON(c, http.StatusOK, m, err)

	m = model.Manufacturing{}
	// calculate the manufacturing for the builder
	if err = manufacturing.NewManufacturing(nil, int32(typeID), 10, 20, 0.1, &m); err == nil {
		//cache.WriteCachedObject(m)
		db.UpdateProfit(m)
	}
}

func GetManufacturingProducts(c *gin.Context) {
	//character := r.Context().Value(CharacterContext).(*model.Character)

	array := strings.Split(c.Query(QueryParamCategoryIDs), SeparatorCategoryIDs)

	categoryIDs := map[int]bool{}

	for _, v := range array {
		i, _ := strconv.Atoi(v)
		categoryIDs[i] = true
	}

	options := db.NewSearchOptions()

	options.NameFilter = c.Query(QueryParamNameFilter)
	options.CategoryIDs = categoryIDs
	options.MaxProductionCosts, _ = httputil.FloatQuery(c, QueryParamMaxProductionCosts)
	options.HasRequiredSkillsOnly, _ = strconv.ParseBool(c.Query(QueryParamHasRequiredSkillsOnly))

	if sortBy := c.Query(QueryParamSortBy); sortBy != "" {
		array = strings.Split(sortBy, SeparatorSortBy)

		options.SortByField = array[0]

		if len(array) > 1 {
			options.SortByDirection = array[1]
		}
	}

	//types, err := cache.GetProductTypes(options, *character)
	types, err := db.GetProductTypes(options)

	httputil.JSON(c, http.StatusOK, types, err)
}
