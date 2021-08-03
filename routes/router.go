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
	"reflect"
	"strconv"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
	"github.com/oxisto/titan/routes/auth"

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
	options.JWTClaims = &model.APIClaims{}
	options.TokenExtractor = auth.ExtractFromFirstAvailable(
		auth.ExtractTokenFromCookie("token"),
		auth.ExtractTokenFromHeader)

	handler := auth.NewHandler(options)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist/titan-frontend", false)))

	r.GET("/auth/login", Login)
	r.GET("/auth/callback", Callback)

	api := r.Group("/api")
	api.Use(handler.AuthRequired)
	api.Use(CharacterRequired)
	{
		character := api.Group("/character")
		{
			character.GET("", GetCharacter)
		}

		corporation := api.Group("/corporation")
		{
			corporation.GET("", GetCorporation)
			corporation.GET("wallets", GetCorporationWallets)
		}

		manufacturing := api.Group("/manufacturing")
		{
			manufacturing.GET("", GetManufacturingProducts)
			manufacturing.GET(":id", GetManufacturing)
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
	claims, ok := c.Value(auth.ClaimsContext).(*model.APIClaims)
	if !ok {
		c.Abort()
		return
	}

	character := &model.Character{}
	cache.GetCharacter(int32(claims.CharacterID), character)

	if limitToCorporationId != 0 && character.CorporationID != limitToCorporationId {
		c.String(http.StatusForbidden, "The corporation you are in is not allowed to access this service")
		c.Abort()
		return
	}

	c.Set(CharacterContext, character)
	c.Next()
}

func IntParam(c *gin.Context, key string) (i int64, err error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

func FloatParam(c *gin.Context, key string) (i float64, err error) {
	return strconv.ParseFloat(c.Param(key), 64)
}

func IntQuery(c *gin.Context, key string) (i int64, err error) {
	return strconv.ParseInt(c.Query(key), 10, 64)
}

func FloatQuery(c *gin.Context, key string) (i float64, err error) {
	return strconv.ParseFloat(c.Query(key), 64)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(c *gin.Context, status int, value interface{}, err error) {
	if err != nil {
		logrus.Errorf("An error occurred during processing of a REST request: %s", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		c.JSON(http.StatusNotFound, nil)
	} else {
		c.JSON(status, value)
	}
}
