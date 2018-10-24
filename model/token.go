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
	"fmt"
	"time"
)

type AccessToken struct {
	expireDate    *time.Time
	Token         string
	CharacterID   int32
	CharacterName string
}

func (token *AccessToken) ID() int32 {
	return token.CharacterID
}

func (token *AccessToken) HashKey() string {
	return fmt.Sprintf("accesstoken:%d", token.ID())
}

func (token *AccessToken) ExpiresOn() *time.Time {
	return token.expireDate
}

func (token *AccessToken) SetExpire(t *time.Time) {
	token.expireDate = t
}

type RefreshToken struct {
	Token       string
	CharacterID int32
}

func (token *RefreshToken) ID() int32 {
	return token.CharacterID
}

func (token *RefreshToken) HashKey() string {
	return fmt.Sprintf("refreshtoken:%d", token.ID())
}

func (token *RefreshToken) ExpiresOn() *time.Time {
	// never expires
	return nil
}

func (token *RefreshToken) SetExpire(t *time.Time) {
	// ignore
}
