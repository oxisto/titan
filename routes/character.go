package routes

import (
	"net/http"

	"github.com/oxisto/titan/model"
)

func GetCharacter(w http.ResponseWriter, r *http.Request) {
	character := r.Context().Value(CharacterContext).(*model.Character)

	JsonResponse(w, r, character, nil)
}
