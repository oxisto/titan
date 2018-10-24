package cache

import (
	sso "github.com/oxisto/evesso"
)

var SSO sso.SingleSignOn

func InitSSO(clientID string, secretKey string, redirectURI string) bool {
	if clientID == "" || secretKey == "" || redirectURI == "" {
		return false
	}

	SSO = sso.SingleSignOn{
		ClientID:    clientID,
		SecretKey:   secretKey,
		RedirectURI: redirectURI,
		Server:      sso.LIVE_SERVER,
	}

	return true
}
