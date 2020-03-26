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

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/oxisto/go-httputil/auth"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"

	"github.com/sirupsen/logrus"
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

func NewRouter(corporationId int32) *gin.Engine {
	limitToCorporationId = corporationId

	options := auth.DefaultOptions
	options.JWTKeySupplier = func(token *jwt.Token) (interface{}, error) {
		return []byte(model.JwtSecretKey), nil
	}
	options.TokenExtractor = auth.ExtractFromFirstAvailable(
		auth.ExtractTokenFromCookie("auth"),
		auth.ExtractTokenFromHeader)

	handler := auth.NewHandler(options)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))

	r.POST("/auth/login", Login)
	r.POST("/auth/callback", Callback)

	api := r.Group("/api")
	api.Use(handler.AuthRequired)
	api.Use(CharacterRequired)
	{
		character := api.Group("/character")
		{
			character.GET("/", GetCharacter)
		}

		corporation := api.Group("/corporation")
		{
			corporation.GET("/", GetCorporation)
		}

		manufacturing := api.Group("/manufacturing")
		{
			manufacturing.GET("/", GetManufacturingProducts)
			manufacturing.GET("/:id", GetManufacturing)
		}
		api.GET("/manufacturing-categories", GetManufacturingCategories)

		industry := api.Group("/industry")
		{
			industry.GET("/jobs", GetIndustryJobs)
		}

		market := api.Group("/market")
		{
			market.POST("/:view", OpenMarketDetail)
		}
	}

	return r
}

func CharacterRequired(c *gin.Context) {
	token, ok := c.Value(auth.ClaimsContext).(*jwt.Token)
	if !ok {
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		c.Abort()
		return
	}

	characterID, ok := (*claims)["CharacterID"].(float64)
	if !ok {
		c.Abort()
		return
	}

	character := &model.Character{}
	cache.GetCharacter(int32(characterID), character)

	if limitToCorporationId != 0 && character.CorporationID != limitToCorporationId {
		c.String(http.StatusForbidden, "The corporation you are in is not allowed to access this service")
		c.Abort()
		return
	}

	c.Set(CharacterContext, character)
	c.Next()
}
