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

package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/manufacturing"
	"github.com/oxisto/titan/model"
)

const (
	QueryParamCategoryIDs           = "categoryIDs"
	QueryParamNameFilter            = "nameFilter"
	QueryParamSortBy                = "sortBy"
	QueryParamMaxProductionCosts    = "maxProductionCosts"
	QueryParamHasRequiredSkillsOnly = "hasRequiredSkillsOnly"
	QueryParamME = "ME"
	QueryParamTE = "TE"

	RouteVarsTypeID = "typeID"

	SeparatorCategoryIDs = ","
	SeparatorSortBy      = ":"
)

func GetManufacturingCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := db.GetCategories()

	JsonResponse(w, r, categories, err)
}

type ManufacturingResponse struct {
	TypeID int32
	Type model.Type
	Manufacturing *manufacturing.Manufacturing
}

func GetManufacturing(w http.ResponseWriter, r *http.Request) {
	var (
		typeID int
		err    error
	)

	ME, err := strconv.Atoi(r.URL.Query().Get(QueryParamME))
	TE, err := strconv.Atoi(r.URL.Query().Get(QueryParamTE))

	character := r.Context().Value(CharacterContext).(*model.Character)

	if typeID, err = strconv.Atoi(mux.Vars(r)["typeID"]); err != nil {
		JsonResponse(w, r, nil, err)
		return
	}

	log.Debugf("Calculating manufacturing information for typeID %d...", typeID)

	resp := ManufacturingResponse{}	
	m := manufacturing.Manufacturing{}

	// calculate it fresh
	if err = manufacturing.NewManufacturing(character, int32(typeID), ME, TE, &m); err != nil {
		// ignore err, but the product should be set
		err = nil
		resp.Type = m.Product
		// keep manufacturing to nil
	} else {
		resp.Manufacturing = &m
	}

	JsonResponse(w, r, m, err)	

	m = manufacturing.Manufacturing{}
	// calculate the manufacturing for the builder
	if err = manufacturing.NewManufacturing(nil, int32(typeID), 10, 20, &m); err == nil {
		cache.WriteCachedObject(m)
	}
}

func GetManufacturingProducts(w http.ResponseWriter, r *http.Request) {
	character := r.Context().Value(CharacterContext).(*model.Character)

	array := strings.Split(r.URL.Query().Get(QueryParamCategoryIDs), SeparatorCategoryIDs)

	categoryIDs := map[int]bool{}

	for _, v := range array {
		i, _ := strconv.Atoi(v)
		categoryIDs[i] = true
	}

	options := cache.NewSearchOptions()

	options.NameFilter = r.URL.Query().Get(QueryParamNameFilter)
	options.CategoryIDs = categoryIDs
	options.MaxProductionCosts, _ = strconv.ParseFloat(r.URL.Query().Get(QueryParamMaxProductionCosts), 64)
	options.HasRequiredSkillsOnly, _ = strconv.ParseBool(r.URL.Query().Get(QueryParamHasRequiredSkillsOnly))

	if sortBy := r.URL.Query().Get(QueryParamSortBy); sortBy != "" {
		array = strings.Split(sortBy, SeparatorSortBy)

		options.SortByField = array[0]

		if len(array) > 1 {
			options.SortByDirection = array[1]
		}
	}

	types, err := cache.GetProductTypes(options, *character)

	JsonResponse(w, r, types, err)
}
