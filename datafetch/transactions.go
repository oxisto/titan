package datafetch

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

type transactionFetcher struct {
	corporationID int32
	division      int32

	log *logrus.Entry
}

func NewTransactionFetcher(corporationID int32, division int32) DataFetcher {
	return &transactionFetcher{
		corporationID: corporationID,
		division:      division,
		log: log.WithFields(logrus.Fields{
			"data":          "transactions",
			"corporationID": corporationID,
			"division":      division,
		}),
	}
}

func (f transactionFetcher) StartLoop() {
	for {
		f.log.Printf("Fetching transactions...")
		duration, err := f.Fetch()

		if err != nil {
			f.log.Printf("An error occured while fetching transactions: %v", err)
		}

		if duration < 0 {
			duration = f.MaxCacheTime()
		}

		f.log.Printf("Waiting for %.2f minutes until next fetch", duration.Minutes())

		time.Sleep(duration)
	}
}

func (f transactionFetcher) MaxCacheTime() time.Duration {
	return time.Duration(3600 * 1000 * 1000)
}

func (f transactionFetcher) Fetch() (time.Duration, error) {
	// find access token for corporation
	accessToken := model.AccessToken{}
	err := cache.GetAccessTokenForCorporation(f.corporationID, &accessToken)
	if err != nil {
		return f.MaxCacheTime(), err
	}

	var options esi.GetCorporationsCorporationIdWalletsDivisionTransactionsOpts

	response, httpResponse, err := cache.ESI.WalletApi.GetCorporationsCorporationIdWalletsDivisionTransactions(
		context.WithValue(context.Background(),
			goesi.ContextAccessToken,
			accessToken.Token),
		f.corporationID,
		f.division,
		&options)
	if err != nil {
		return f.MaxCacheTime(), err
	}

	t, err := time.Parse(time.RFC1123, httpResponse.Header.Get("Expires"))
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	// loop through all transactions
	for _, t := range response {
		transaction := model.Transaction{
			TransactionID: t.TransactionId,
			CorporationID: int64(f.corporationID),
			Division:      f.division,
			ClientID:      t.ClientId,
			Date:          t.Date,
			IsBuy:         t.IsBuy,
			JournalRefID:  t.JournalRefId,
			LocationID:    t.LocationId,
			Quantity:      int(t.Quantity),
			TypeID:        model.TypeIdentifier(t.TypeId),
			UnitPrice:     t.UnitPrice,
		}

		f.log.Printf("Discovered new transaction %d (%d, %d x %.2f ISK))", transaction.TransactionID, transaction.TypeID, transaction.Quantity, transaction.UnitPrice)

		if err := db.InsertTransaction(&transaction); err != nil {
			f.log.Printf("Could not insert transaction with ID %d: %v", transaction.TransactionID, err.Error())
		}
	}

	return time.Until(t), nil
}
