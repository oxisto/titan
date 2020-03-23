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

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
)

func Login(w http.ResponseWriter, r *http.Request) {
	scope := "publicData esi-skills.read_skills.v1 esi-corporations.read_corporation_membership.v1 esi-ui.open_window.v1 esi-wallet.read_corporation_wallets.v1 esi-assets.read_corporation_assets.v1 esi-corporations.read_blueprints.v1 esi-industry.read_corporation_jobs.v1"

	t := time.Now()
	state := base64.StdEncoding.EncodeToString([]byte(t.String()))

	w.Header().Add("Location", cache.SSO.Redirect(state, &scope))
	w.WriteHeader(http.StatusFound)
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Errorf("Could not fetch access token: %v", err)
	w.Header().Add("Location", "/auth/login")
	w.WriteHeader(http.StatusFound)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// fetch access token with authorization code
	tokenResponse, expiryTime, characterID, characterName, err := cache.SSO.AccessToken(code, false)
	if err != nil {
		HandleError(err, w, r)
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
		HandleError(err, w, r)
		return
	}

	// cache the refresh token (they never expire)
	err = cache.WriteCachedObject(&model.RefreshToken{
		CharacterID: int32(characterID),
		Token:       tokenResponse.RefreshToken,
	})
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// issue an authentication token for our own API
	authToken, err := model.IssueToken(int32(characterID))
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// redirect to main dashboard page
	w.Header().Add("Location", "/#?token="+authToken)
	w.Header().Add("Set-Cookie", "token="+authToken+"; Path=/")
	w.WriteHeader(http.StatusFound)
}
