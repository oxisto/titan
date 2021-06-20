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
	JobID                int32      `json:"jobID"`
	ActivityID           int32      `json:"activityID"`
	CompletedCharacterID int32      `json:"completedCharacterID"`
	CompletedDate        *time.Time `json:"completedDate"`
	Cost                 float64    `json:"cost"`
	Duration             int32      `json:"duration"`
	EndDate              time.Time  `json:"endDate"`
	FacilityID           int64      `json:"facilityID"`
	InstallerID          int32      `json:"installerID"`
	LocationID           int64      `json:"locationID"`
	BlueprintID          int64      `json:"blueprintID"`
	BlueprintTypeID      int32      `json:"blueprintTypeID"`
	StartDate            time.Time  `json:"startDate"`
	PauseDate            *time.Time `json:"pausedDate"`
	LicensedRuns         int32      `json:"licensedRuns"`
	OutputLocationID     int64      `json:"outputLocationID"`
	Probability          float32    `json:"probability"`
	ProductTypeID        int32      `json:"productTypeID"`
	Runs                 int32      `json:"runs"`
	SuccesfulRuns        int32      `json:"succesfulRuns"`
	Status               string     `json:"status"`
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
