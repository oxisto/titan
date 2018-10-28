/*
Copyright 2018 Christian Banse

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
	"time"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
)

func Login(w http.ResponseWriter, r *http.Request) {
	scope := "publicData esi-skills.read_skills.v1 esi-corporations.read_corporation_membership.v1 esi-ui.open_window.v1 esi-wallet.read_corporation_wallets.v1"
	w.Header().Add("Location", cache.SSO.Redirect(nil, &scope))
	w.WriteHeader(http.StatusFound)
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Errorf("Could not fetch access token: %v", err)
	w.Header().Add("Location", "/auth/login")
	w.WriteHeader(http.StatusFound)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	// fetch access token with authorization code
	tokenResponse, err := cache.SSO.AccessToken(r.URL.Query().Get("code"), false)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// verify access token
	verifyResponse, err := cache.SSO.Verify(tokenResponse.AccessToken)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// create a new access token cache object
	accessToken := model.AccessToken{
		CharacterID:   verifyResponse.CharacterID,
		CharacterName: verifyResponse.CharacterName,
		Token:         tokenResponse.AccessToken,
	}

	t, err := time.Parse(time.RFC3339, verifyResponse.ExpiresOn+"Z")
	if err != nil {
		HandleError(err, w, r)
		return
	}
	accessToken.SetExpire(&t)

	// cache the access token
	err = cache.WriteCachedObject(&accessToken)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// cache the refresh token (they never expire)
	err = cache.WriteCachedObject(&model.RefreshToken{
		CharacterID: verifyResponse.CharacterID,
		Token:       tokenResponse.RefreshToken,
	})
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// issue an authentication token for our own API
	authToken, err := model.IssueToken(verifyResponse.CharacterID)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	// redirect to main dashboard page
	w.Header().Add("Location", "/#?token="+authToken)
	w.Header().Add("Set-Cookie", "token="+authToken+"; Path=/")
	w.WriteHeader(http.StatusFound)
}
