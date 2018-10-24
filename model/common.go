package model

import (
	"time"
)

type CachedObject interface {
	ID() int32
	HashKey() string
	ExpiresOn() *time.Time
	SetExpire(t *time.Time)
}

func SafeInt32(obj *int32) int32 {
	if obj == nil {
		return int32(0)
	} else {
		return int32(*obj)
	}
}

func SafeInt64(obj *int64) int64 {
	if obj == nil {
		return int64(0)
	} else {
		return int64(*obj)
	}
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}
