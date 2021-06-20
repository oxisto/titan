package datafetch

import (
	"context"
	"net/http"
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
	metadata
}

func NewIndustryJobsFetcher() DataFetcher {
	return &industryJobsFetcher{
		metadata: metadata{
			dataType:     "industry-jobs",
			maxCacheTime: time.Minute * 5,
		},
	}
}

func (i *industryJobsFetcher) Fetch(ctx FetchContext) (*http.Response, error) {
	var options esi.GetCorporationsCorporationIdIndustryJobsOpts
	options.IncludeCompleted = optional.NewBool(true)

	if ctx.lastETag != "" {
		options.IfNoneMatch = optional.NewString(ctx.lastETag)
	}

	response, httpResponse, err := cache.ESI.IndustryApi.GetCorporationsCorporationIdIndustryJobs(
		context.WithValue(context.Background(),
			goesi.ContextAccessToken,
			ctx.accessToken.Token),
		ctx.corporationID,
		&options)
	if err != nil {
		return httpResponse, err
	}

	limitFields := logrus.Fields{
		"esi-err-remain": httpResponse.Header.Get("x-esi-error-limit-remain"),
		"esi-err-reset":  httpResponse.Header.Get("x-esi-error-limit-reset"),
	}

	if httpResponse.StatusCode != 304 {
		ctx.log.WithFields(limitFields).Infof("Retrieved %d industry jobs", len(response))

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

			ctx.log.Debugf("Discovered industry job %d (%d, %d)", job.JobID, job.ActivityID, job.BlueprintTypeID)

			if err := db.UpdateIndustryJob(&job); err != nil {
				ctx.log.Errorf("Could not update industry job ID %d: %v", job.JobID, err)
			}
		}
	} else {
		ctx.log.WithFields(limitFields).Info("Industry jobs have not changed")
	}

	return httpResponse, nil
}
