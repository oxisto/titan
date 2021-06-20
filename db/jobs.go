package db

import "github.com/oxisto/titan/model"

func GetIndustryJobs(corporationID int32) ([]*model.IndustryJobs, error) {
	return nil, nil
}

func UpdateIndustryJob(job *model.IndustryJob) error {
	_, err := pdb.Exec(`INSERT INTO
	"industryJobs"
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	ON CONFLICT("jobID") DO UPDATE
	SET
		"activityID" = $2,
		"completedCharacterID" = $3,
		"completedDate" = $4,
		"cost" = $5,
		"duration" = $6,
		"endDate" = $7,
		"facilityID" = $8,
		"installerID" = $9,
		"locationID" = $10,
		"blueprintID" = $11,
		"blueprintTypeID" = $12,
		"startDate" = $13,
		"pauseDate" = $14,
		"licensedRuns" = $15,
		"outputLocationID" = $16,
		"probability" = $17,
		"productTypeID" = $18,
		"runs" = $19,
		"succesfulRuns" = $20,
		"status" = $21`,
		job.JobID,
		job.ActivityID,
		job.CompletedCharacterID,
		job.CompletedDate,
		job.Cost,
		job.Duration,
		job.EndDate,
		job.FacilityID,
		job.InstallerID,
		job.LocationID,
		job.BlueprintID,
		job.BlueprintTypeID,
		job.StartDate,
		job.PauseDate,
		job.LicensedRuns,
		job.OutputLocationID,
		job.Probability,
		job.ProductTypeID,
		job.Runs,
		job.SuccesfulRuns,
		job.Status)

	return err
}
