package datafetch

import (
	"context"
	"time"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
)

type industryJobsFetcher struct {
	corporationID int32

	log      *logrus.Entry
	lastETag string
}

func NewIndustryJobsFetcher(corporationID int32) DataFetcher {
	return &industryJobsFetcher{
		corporationID: corporationID,
		log: log.WithFields(logrus.Fields{
			"data":          "industry-jobs",
			"corporationID": corporationID,
		}),
	}
}

func (i industryJobsFetcher) StartLoop() {
	for {
		i.log.Printf("Fetching industry jobs...")
		duration, err := i.Fetch()

		if err != nil {
			i.log.Printf("An error occured while fetching jobs: %v", err)
		}

		if duration < 0 {
			duration = i.MaxCacheTime()
		}

		i.log.Printf("Waiting for %.2f minutes until next fetch", duration.Minutes())

		time.Sleep(duration)
	}
}

func (i industryJobsFetcher) MaxCacheTime() time.Duration {
	return time.Minute * 5
}

func (i *industryJobsFetcher) Fetch() (time.Duration, error) {
	// find access token for corporation
	accessToken := model.AccessToken{}
	err := cache.GetAccessTokenForCorporation(i.corporationID, &accessToken)
	if err != nil {
		return i.MaxCacheTime(), err
	}

	var options esi.GetCorporationsCorporationIdIndustryJobsOpts
	options.IncludeCompleted = optional.NewBool(true)

	if i.lastETag != "" {
		options.IfNoneMatch = optional.NewString(i.lastETag)
	}

	response, httpResponse, err := cache.ESI.IndustryApi.GetCorporationsCorporationIdIndustryJobs(
		context.WithValue(context.Background(),
			goesi.ContextAccessToken,
			accessToken.Token),
		i.corporationID,
		&options)
	if err != nil {
		return i.MaxCacheTime(), err
	}

	t, err := time.Parse(time.RFC1123, httpResponse.Header.Get("Expires"))
	if err != nil {
		return i.MaxCacheTime(), err
	}

	limitFields := logrus.Fields{
		"esi-err-remain": httpResponse.Header.Get("x-esi-error-limit-remain"),
		"esi-err-reset":  httpResponse.Header.Get("x-esi-error-limit-reset"),
	}

	if httpResponse.StatusCode != 304 {
		i.log.WithFields(limitFields).Infof("Retrieved %d industry jobs", len(response))

		// store the ETag
		i.lastETag = httpResponse.Header.Get("etag")

		// loop through all jobs
		for _, t := range response {
			job := model.IndustryJob{
				ActivityID:           t.ActivityId,
				BlueprintID:          t.BlueprintId,
				BlueprintTypeID:      t.BlueprintTypeId,
				CompletedCharacterID: t.CompletedCharacterId,
				Cost:                 t.Cost,
				Duration:             t.Duration,
				EndDate:              t.EndDate,
				FacilityID:           t.FacilityId,
				InstallerID:          t.InstallerId,
				LicensedRuns:         t.LicensedRuns,
				LocationID:           t.LocationId,
				OutputLocationID:     t.OutputLocationId,
				Probability:          t.Probability,
				ProductTypeID:        t.ProductTypeId,
				Runs:                 t.Runs,
				JobID:                t.JobId,
				StartDate:            t.StartDate,
				Status:               t.Status,
				SuccesfulRuns:        t.SuccessfulRuns,
			}

			if !t.CompletedDate.IsZero() {
				job.CompletedDate = &t.CompletedDate
			}

			if !t.PauseDate.IsZero() {
				job.PauseDate = &t.PauseDate
			}

			i.log.Debugf("Discovered industry job %d (%d, %d)", job.JobID, job.ActivityID, job.BlueprintTypeID)

			if err := db.UpdateIndustryJob(&job); err != nil {
				i.log.Printf("Could not update industry job ID %d: %v", job.JobID, err)
			}
		}
	} else {
		i.log.WithFields(limitFields).Info("Industry jobs have not changed")
	}

	return time.Until(t), nil
}
