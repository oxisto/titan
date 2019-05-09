package finance

import (
	"context"
	"time"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "finance")
}

func FetchJournal(corporationID int32, journal int32) (time.Duration, error) {
	// find access token for corporation
	accessToken := model.AccessToken{}
	err := cache.GetAccessTokenForCorporation(corporationID, &accessToken)
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	// for now, no paging
	options := esi.GetCorporationsCorporationIdWalletsDivisionJournalOpts{
		//Page: optional.NewInt32(1),
	}

	response, httpResponse, err := cache.ESI.WalletApi.GetCorporationsCorporationIdWalletsDivisionJournal(context.WithValue(context.Background(), goesi.ContextAccessToken, accessToken.Token), corporationID, 1, &options)
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	t, err := time.Parse(time.RFC1123, httpResponse.Header.Get("Expires"))
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	entryIDs, err := db.GetJournalEntryIDs()
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	// convert to map for easier searching
	entryMap := map[int64]bool{}
	for _, ID := range entryIDs {
		entryMap[ID] = true
	}

	// loop through all journal entries
	for _, journal := range response {
		// skip, if it is already known to us
		if entryMap[journal.Id] {
			continue
		}

		entry := model.JournalEntry{
			Amount:        journal.Amount,
			Balance:       journal.Balance,
			Date:          journal.Date,
			Description:   journal.Description,
			FirstPartyID:  journal.FirstPartyId,
			ID:            journal.Id,
			RefType:       journal.RefType,
			SecondPartyID: journal.SecondPartyId,
		}

		log.Printf("Discovered new journal entry %d (%s, %.2f ISK))", entry.ID, entry.Description, entry.Amount)

		if err := db.InsertJournalEntry(entry); err != nil {
			log.Printf("Could not insert journal entry with ID %d: %v", entry.ID, err.Error())
		}
	}

	return t.Sub(time.Now()), nil
}
