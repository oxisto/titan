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
