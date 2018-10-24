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

	esi "github.com/evecentral/esiapi/client"
	esiUI "github.com/evecentral/esiapi/client/user_interface"
	"github.com/go-openapi/runtime/client"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
)

func OpenMarketDetail(w http.ResponseWriter, r *http.Request) {
	character := r.Context().Value(CharacterContext).(*model.Character)

	if typeID, err := strconv.Atoi(r.URL.Query().Get("typeID")); err == nil {
		OpenMarket(character.ID(), int32(typeID), w, r)
	}
}

func OpenMarket(characterID int32, typeID int32, w http.ResponseWriter, r *http.Request) {
	uiParams := esiUI.NewPostUIOpenwindowMarketdetailsParams()
	uiParams.TypeID = typeID

	// find access token for character
	accessToken := model.AccessToken{}
	err := cache.GetAccessToken(characterID, &accessToken)
	if err != nil {
		JsonResponse(w, r, nil, err)
		return
	}

	_, err = esi.Default.UserInterface.PostUIOpenwindowMarketdetails(uiParams, client.BearerToken(accessToken.Token))
	if err != nil {
		JsonResponse(w, r, nil, err)
		return
	}
}
