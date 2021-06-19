package model

import (
	"fmt"
	"time"
)

type SystemCostIndex struct {
	expireDate    *time.Time
	ActivityCost  map[string]float32
	SolarSystemID int32
}

func (i *SystemCostIndex) ID() int32 {
	return i.SolarSystemID
}

func (i *SystemCostIndex) ExpiresOn() *time.Time {
	return i.expireDate
}

func (i *SystemCostIndex) SetExpire(t *time.Time) {
	i.expireDate = t
}

func (i *SystemCostIndex) HashKey() string {
	return fmt.Sprintf("system-cost-index:%d", i.ID())
}

type IndustryJobs struct {
	expireDate    *time.Time
	CorporationID int32                  `json:"corporationID"`
	Jobs          map[string]IndustryJob `json:"jobs"`
}

type IndustryJob struct {
	ActivityID       int32   `json:"activityID"`
	Blueprint        *Type   `json:"blueprint"`
	StartDate        int64   `json:"startDate"`
	EndDate          int64   `json:"endDate"`
	CompletedDate    int64   `json:"completedDate"`
	PauseDate        int64   `json:"pausedDate"`
	LicensedRuns     int     `json:"licensedRuns"`
	OutputLocationID int64   `json:"outputLocationID"`
	Probability      float32 `json:"probability"`
	SuccesfulRuns    int     `json:"succesfulRuns"`
	Status           string  `json:"status"`
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
