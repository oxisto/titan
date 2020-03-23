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
	"context"
	"net/http"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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
	router.Handle("/api/character", WithMiddleware(middleware, GetCharacter))
	router.Handle("/api/corporation", WithMiddleware(middleware, GetCorporation))
	router.Handle("/api/manufacturing", WithMiddleware(middleware, GetManufacturingProducts))
	router.Handle("/api/manufacturing/{"+RouteVarsTypeID+"}", WithMiddleware(middleware, GetManufacturing))
	router.Handle("/api/manufacturing-categories", WithMiddleware(middleware, GetManufacturingCategories))
	router.Handle("/api/industry/jobs", WithMiddleware(middleware, GetIndustryJobs))
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
		http.Error(w, "The corporation you are in is not allowed to access this service", http.StatusForbidden)
		return
	}

	request := r.WithContext(context.WithValue(r.Context(), CharacterContext, character))

	*r = *request
	next(w, r)
}
