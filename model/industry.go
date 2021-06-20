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
	CorporationID int32                       `json:"corporationID"`
	Jobs          []*IndustryJobWithTypeNames `json:"jobs"`
}

type IndustryJob struct {
	JobID                int32      `json:"jobID" db:"jobID"`
	ActivityID           int32      `json:"activityID" db:"activityID"`
	CompletedCharacterID int32      `json:"completedCharacterID" db:"completedCharacterID"`
	CompletedDate        *time.Time `json:"completedDate" db:"completedDate"`
	Cost                 float64    `json:"cost" db:"cost"`
	Duration             int32      `json:"duration" db:"duration"`
	EndDate              time.Time  `json:"endDate" db:"endDate"`
	FacilityID           int64      `json:"facilityID" db:"facilityID"`
	InstallerID          int32      `json:"installerID" db:"installerID"`
	LocationID           int64      `json:"locationID" db:"locationID"`
	BlueprintID          int64      `json:"blueprintID" db:"blueprintID"`
	BlueprintTypeID      int32      `json:"blueprintTypeID" db:"blueprintTypeID"`
	StartDate            time.Time  `json:"startDate" db:"startDate"`
	PauseDate            *time.Time `json:"pauseDate" db:"pauseDate"`
	LicensedRuns         int32      `json:"licensedRuns" db:"licensedRuns"`
	OutputLocationID     int64      `json:"outputLocationID" db:"outputLocationID"`
	Probability          float32    `json:"probability" db:"probability"`
	ProductTypeID        int32      `json:"productTypeID" db:"productTypeID"`
	Runs                 int32      `json:"runs" db:"runs"`
	SuccesfulRuns        int32      `json:"succesfulRuns" db:"succesfulRuns"`
	Status               string     `json:"status" db:"status"`
}

type IndustryJobWithTypeNames struct {
	*IndustryJob
	BlueprintTypeName string `json:"blueprintTypeName" db:"blueprintTypeName"`
}
