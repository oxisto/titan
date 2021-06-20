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

type journalFechter struct {
	metadata

	division int32
}

func NewJournalFetcher(division int32) DataFetcher {
	return &journalFechter{
		division: division,
		metadata: metadata{
			dataType:         "journal",
			maxCacheTime:     time.Hour,
			additionalFields: logrus.Fields{"division": division},
		},
	}
}

func (f *journalFechter) Fetch(ctx FetchContext) (*http.Response, error) {
	var options esi.GetCorporationsCorporationIdWalletsDivisionJournalOpts
	if ctx.lastETag != "" {
		options.IfNoneMatch = optional.NewString(ctx.lastETag)
	}

	response, httpResponse, err := cache.ESI.WalletApi.GetCorporationsCorporationIdWalletsDivisionJournal(
		context.WithValue(context.Background(),
			goesi.ContextAccessToken,
			ctx.accessToken.Token),
		ctx.corporationID,
		f.division,
		&options)
	if err != nil {
		return httpResponse, err
	}

	limitFields := logrus.Fields{
		"esi-err-remain": httpResponse.Header.Get("x-esi-error-limit-remain"),
		"esi-err-reset":  httpResponse.Header.Get("x-esi-error-limit-reset"),
	}

	if httpResponse.StatusCode != 304 {
		ctx.log.WithFields(limitFields).Infof("Retrieved %d journal entries", len(response))

		// loop through all journal entries
		for _, journal := range response {
			entry := model.JournalEntry{
				Amount:        journal.Amount,
				Balance:       journal.Balance,
				Date:          journal.Date,
				Description:   journal.Description,
				FirstPartyID:  journal.FirstPartyId,
				ID:            journal.Id,
				RefType:       journal.RefType,
				SecondPartyID: journal.SecondPartyId,
				CorporationID: ctx.corporationID,
				Division:      1,
			}

			ctx.log.Debugf("Discovered new journal entry %d (%s, %.2f ISK))", entry.ID, entry.Description, entry.Amount)

			if err := db.InsertJournalEntry(entry); err != nil {
				ctx.log.Errorf("Could not insert journal entry with ID %d: %v", entry.ID, err.Error())
			}
		}
	} else {
		ctx.log.WithFields(limitFields).Info("Journal has not changed")
	}

	return httpResponse, nil
}
