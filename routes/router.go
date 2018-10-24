package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	CharacterContext = "character"
)

var (
	limitToCorporationId int32
	log                  *logrus.Entry
)

func init() {
	log = logrus.WithField("component", "routes")
}

func JsonResponse(w http.ResponseWriter, r *http.Request, object interface{}, err error) {
	// uh-uh, we have an error
	if err != nil {
		log.Error("An error occured during processing of a REST request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return not found if object is nil
	if object == nil {
		http.NotFound(w, r)
		return
	}

	// otherwise, lets try to decode the JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(object); err != nil {
		// uh-uh we couldn't decode the JSON
		log.Errorf("An error occured during encoding of the JSON response: %v. Payload was: %+v", err, object)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewRouter(corporationId int32) *mux.Router {
	limitToCorporationId = corporationId

	middleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(model.JwtSecretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth/callback", Callback)
	router.HandleFunc("/auth/login", Login)
	router.HandleFunc("/slack/callback", SlackCallback)
	router.Handle("/api/character", WithMiddleware(middleware, GetCharacter))
	router.Handle("/api/corporation", WithMiddleware(middleware, GetCorporation))
	router.Handle("/api/manufacturing", WithMiddleware(middleware, GetManufacturingProducts))
	router.Handle("/api/manufacturing/{"+RouteVarsTypeID+"}", WithMiddleware(middleware, GetManufacturing))
	router.Handle("/api/manufacturing-categories", WithMiddleware(middleware, GetManufacturingCategories))
	router.Handle("/api/market/view", WithMiddleware(middleware, OpenMarketDetail))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist")))

	return router
}

func WithMiddleware(middleware *jwtmiddleware.JWTMiddleware, handlerFunc http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(middleware.HandlerWithNext),
		negroni.HandlerFunc(HandleFetchCharacterWithNext),
		negroni.Wrap(handlerFunc),
	)
}

// Special implementation for Negroni, but could be used elsewhere.
func HandleFetchCharacterWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, ok := r.Context().Value("user").(*jwt.Token)
	if !ok {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	characterID, ok := claims["CharacterID"].(float64)
	if !ok {
		return
	}

	character := &model.Character{}
	cache.GetCharacter(int32(characterID), character)

	if limitToCorporationId != 0 && character.CorporationID != limitToCorporationId {
		return
	}

	request := r.WithContext(context.WithValue(r.Context(), CharacterContext, character))

	*r = *request
	next(w, r)
}
