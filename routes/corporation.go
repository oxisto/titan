package routes

import (
	"net/http"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
)

func GetCorporation(w http.ResponseWriter, r *http.Request) {
	character := r.Context().Value(CharacterContext).(*model.Character)
	corporation := &model.Corporation{}

	err := cache.GetCorporation(character.CharacterID, character.CorporationID, corporation)

	JsonResponse(w, r, corporation, err)
}
