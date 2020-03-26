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
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
)

func Login(c *gin.Context) {
	scope := "publicData esi-skills.read_skills.v1 esi-corporations.read_corporation_membership.v1 esi-ui.open_window.v1 esi-wallet.read_corporation_wallets.v1 esi-assets.read_corporation_assets.v1 esi-corporations.read_blueprints.v1 esi-industry.read_corporation_jobs.v1"

	t := time.Now()
	state := base64.StdEncoding.EncodeToString([]byte(t.String()))

	c.Header("Location", cache.SSO.Redirect(state, &scope))
	c.String(http.StatusFound, "")
}

func HandleError(err error, c *gin.Context) {
	log.Errorf("Could not fetch access token: %v", err)

	c.Header("Location", "/auth/login")
	c.String(http.StatusFound, "")
}

func Callback(c *gin.Context) {
	code := c.Query("code")

	// fetch access token with authorization code
	tokenResponse, expiryTime, characterID, characterName, err := cache.SSO.AccessToken(code, false)
	if err != nil {
		HandleError(err, c)
		return
	}

	// create a new access token cache object
	accessToken := model.AccessToken{
		CharacterID:   int32(characterID),
		CharacterName: characterName,
		Token:         tokenResponse.AccessToken,
	}
	accessToken.SetExpire(&expiryTime)

	// cache the access token
	err = cache.WriteCachedObject(&accessToken)
	if err != nil {
		HandleError(err, c)
		return
	}

	// cache the refresh token (they never expire)
	err = cache.WriteCachedObject(&model.RefreshToken{
		CharacterID: int32(characterID),
		Token:       tokenResponse.RefreshToken,
	})
	if err != nil {
		HandleError(err, c)
		return
	}

	// issue an authentication token for our own API
	authToken, err := model.IssueToken(int32(characterID))
	if err != nil {
		HandleError(err, c)
		return
	}

	// redirect to main dashboard page
	c.Header("Location", "/#?token="+authToken)
	c.Header("Set-Cookie", "token="+authToken+"; Path=/")
	c.String(http.StatusFound, "")
}
