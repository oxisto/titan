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
	QueryParamCategoryIDs        = "categoryIDs"
	QueryParamNameFilter         = "nameFilter"
	QueryParamSortBy             = "sortBy"
	QueryParamMaxProductionCosts = "maxProductionCosts"

	RouteVarsTypeID = "typeID"

	SeparatorCategoryIDs = ","
	SeparatorSortBy      = ":"
)

func GetManufacturingCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := db.GetCategories()

	JsonResponse(w, r, categories, err)
}

func GetManufacturing(w http.ResponseWriter, r *http.Request) {
	var (
		typeID int
		err    error
	)

	character := r.Context().Value(CharacterContext).(*model.Character)

	if typeID, err = strconv.Atoi(mux.Vars(r)["typeID"]); err != nil {
		JsonResponse(w, r, nil, err)
		return
	}

	log.Debugf("Calculating manufacturing information for typeID %d...", typeID)

	m := manufacturing.Manufacturing{}

	// calculate it fresh and update cache
	if err = manufacturing.NewManufacturing(*character, int32(typeID), &m); err == nil {
		cache.WriteCachedObject(m)
	}

	JsonResponse(w, r, m, err)
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
