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
