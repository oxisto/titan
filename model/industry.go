package model

import (
	"fmt"
	"time"
)

type IndustryJobs struct {
	expireDate    *time.Time
	CorporationID int32                  `json:"corporationID"`
	Jobs          map[string]IndustryJob `json:"jobs"`
}

type IndustryJob struct {
	ActivityID int32 `json:"activityID"`
	Blueprint  Type  `json:"blueprint"`
}

func (i *IndustryJobs) ID() int32 {
	return i.CorporationID
}

func (i *IndustryJobs) ExpiresOn() *time.Time {
	return i.expireDate
}

func (i *IndustryJobs) SetExpire(t *time.Time) {
	i.expireDate = t
}

func (i *IndustryJobs) HashKey() string {
	return fmt.Sprintf("industry-jobs:%d", i.ID())
}
