package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	JwtSecretKey = "aCf2enHaMAxKKrNKZgVaMCFn"
)

type APIClaims struct {
	*jwt.StandardClaims
	CharacterID int32
}

// IssueToken issues a JWT token for use of the API
func IssueToken(characterID int32) (token string, err error) {
	key := []byte(JwtSecretKey)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &APIClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
		characterID,
	})

	token, err = claims.SignedString(key)
	return
}
