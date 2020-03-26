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

	"github.com/antihax/goesi"
	"github.com/gin-gonic/gin"
	"github.com/oxisto/titan/model"

	"github.com/oxisto/go-httputil"
	"github.com/oxisto/titan/cache"
)

func OpenMarketDetail(c *gin.Context) {
	character := c.Value(CharacterContext).(*model.Character)

	if typeID, err := httputil.IntParam(c, "id"); err == nil {
		OpenMarket(character.ID(), int32(typeID), c)
	}
}

func OpenMarket(characterID int32, typeID int32, c *gin.Context) {
	// find access token for character
	accessToken := model.AccessToken{}
	err := cache.GetAccessToken(characterID, &accessToken)
	if err != nil {
		httputil.JSON(c, http.StatusNotFound, nil, err)
		return
	}

	_, err = cache.ESI.UserInterfaceApi.PostUiOpenwindowMarketdetails(context.WithValue(context.Background(), goesi.ContextAccessToken, accessToken.Token), typeID, nil)
	if err != nil {
		httputil.JSON(c, http.StatusNotFound, nil, err)
		return
	}
}
