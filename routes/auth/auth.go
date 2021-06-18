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

package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

// TokenExtractorFunc defines a function which extracts a token out of an HTTP request
type TokenExtractorFunc func(r *http.Request) (token string, err error)

// ErrorHandlerFunc defines a function which is called if an error occured during the
// extraction of a token
type ErrorHandlerFunc func(err error, w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type Options struct {
	TokenExtractor TokenExtractorFunc
	ErrorHandler   ErrorHandlerFunc
	JWTKeySupplier jwt.Keyfunc
	JWTClaims      jwt.Claims
	RequireToken   bool
}

type JWTHandler struct {
	options Options
}

var ErrNoToken = errors.New("jwt: no token could be extracted and Options.RequireToken is true")

var DefaultOptions Options

const ClaimsContext string = "claims"

func init() {
	DefaultOptions = Options{
		RequireToken:   true,
		TokenExtractor: ExtractTokenFromHeader,
		ErrorHandler:   nil,
		JWTClaims:      &jwt.StandardClaims{},
	}
}

// NewHandler creates a new instance of the JWT handler
func NewHandler(Options Options) *JWTHandler {
	handler := JWTHandler{Options}

	return &handler
}

// AuthRequired is a middleware for gin-gonic
func (h JWTHandler) AuthRequired(c *gin.Context) {
	err := h.parseJWT(c)

	if err == nil {
		c.Next()
	} else {
		c.Status(http.StatusForbidden)
		c.Abort()
	}
}

func (h JWTHandler) parseJWT(c *gin.Context) (err error) {
	var token string
	var parsed *jwt.Token

	// extract token
	token, err = h.options.TokenExtractor(c.Request)

	if err == nil && token == "" && h.options.RequireToken {
		err = ErrNoToken
		return
	}

	parsed, err = jwt.ParseWithClaims(token, h.options.JWTClaims, h.options.JWTKeySupplier)

	if err != nil {
		return
	}

	// update context
	c.Set(ClaimsContext, parsed.Claims)

	return
}

// ExtractTokenFromHeader extracts a JWT out of the authorization header of an HTTP request
func ExtractTokenFromHeader(r *http.Request) (token string, err error) {
	authorization := strings.Split(r.Header.Get("Authorization"), " ")

	if len(authorization) >= 2 && authorization[0] == "Bearer" {
		return authorization[1], nil
	}

	// no token was found, but also no error occurred
	return "", nil
}

// ExtractTokenFromCookie extracts a JWT out of an HTTP cookie
func ExtractTokenFromCookie(cookie string) TokenExtractorFunc {
	return func(r *http.Request) (token string, err error) {
		cookie, err := r.Cookie(cookie)

		// dont throw error, if cookie is not found, just return empty token
		if err != nil && err == http.ErrNoCookie {
			return "", nil
		}

		if err != nil {
			return "", err
		}

		return cookie.Value, nil
	}
}

// ExtractFromFirstAvailable extracts the token out of the specified extractors.
// The first token that is found will be returned
func ExtractFromFirstAvailable(extractors ...TokenExtractorFunc) TokenExtractorFunc {
	return func(r *http.Request) (token string, err error) {
		for _, extractor := range extractors {
			token, err := extractor(r)

			if err != nil {
				return "", err
			}

			if token != "" {
				return token, nil
			}
		}

		return "", nil
	}
}

// IssueToken is a little helper that issues tokens for a specified key, subject and expiry time
func IssueToken(key []byte, subject string, expiry time.Time) (token *oauth2.Token, err error) {
	var accessToken string

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512,
		&jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
			Subject:   subject,
		},
	)

	if accessToken, err = claims.SignedString(key); err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      expiry,
	}, nil
}
