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
