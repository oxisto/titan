package finance

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

func FetchTransations(corporationID int32, journal int32) (time.Duration, error) {
	// find access token for corporation
	accessToken := model.AccessToken{}
	err := cache.GetAccessTokenForCorporation(corporationID, &accessToken)
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	// get last transactionID
	transaction, err := db.GetLatestTransaction()

	var options esi.GetCorporationsCorporationIdWalletsDivisionTransactionsOpts

	if transaction != nil {
		log.Debugf("Latest transaction was %d", transaction.TransactionID)
		options.FromId = optional.NewInt64(transaction.TransactionID)
	} else {
		log.Debugf("No previous transaction found. Fetching all")
	}

	response, httpResponse, err := cache.ESI.WalletApi.GetCorporationsCorporationIdWalletsDivisionTransactions(context.WithValue(context.Background(), goesi.ContextAccessToken, accessToken.Token), corporationID, 1, &options)
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	t, err := time.Parse(time.RFC1123, httpResponse.Header.Get("Expires"))
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	entryIDs, err := db.GetTransactionIDs()
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	// convert to map for easier searching
	entryMap := map[int64]bool{}
	for _, ID := range entryIDs {
		entryMap[ID] = true
	}

	// loop through all transactions
	for _, t := range response {
		// skip, if it is already known to us
		if entryMap[t.TransactionId] {
			continue
		}

		transaction := model.Transaction{
			TransactionID: t.TransactionId,
			ClientID:      t.ClientId,
			Date:          t.Date,
			IsBuy:         t.IsBuy,
			JournalRefID:  t.JournalRefId,
			LocationID:    t.LocationId,
			Quantity:      int(t.Quantity),
			TypeID:        model.TypeIdentifier(t.TypeId),
			UnitPrice:     t.UnitPrice,
		}

		log.Printf("Discovered new transaction %d (%d, %d x %.2f ISK))", transaction.TransactionID, transaction.TypeID, transaction.Quantity, transaction.UnitPrice)

		if err := db.InsertTransaction(transaction); err != nil {
			log.Printf("Could not insert transaction with ID %d: %v", transaction.TransactionID, err.Error())
		}
	}

	return t.Sub(time.Now()), nil
}
