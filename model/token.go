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
