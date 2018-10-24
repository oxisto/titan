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
