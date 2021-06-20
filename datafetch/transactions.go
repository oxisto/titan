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

type transactionFetcher struct {
	metadata

	division int32
}

func NewTransactionFetcher(division int32) DataFetcher {
	return &transactionFetcher{
		division: division,
		metadata: metadata{
			dataType:         "transactions",
			maxCacheTime:     time.Hour,
			additionalFields: logrus.Fields{"division": division},
		},
	}
}

func (f *transactionFetcher) Fetch(ctx FetchContext) (*http.Response, error) {
	var options esi.GetCorporationsCorporationIdWalletsDivisionTransactionsOpts
	if ctx.lastETag != "" {
		options.IfNoneMatch = optional.NewString(ctx.lastETag)
	}

	response, httpResponse, err := cache.ESI.WalletApi.GetCorporationsCorporationIdWalletsDivisionTransactions(
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
		ctx.log.WithFields(limitFields).Infof("Retrieved %d transactions", len(response))

		// loop through all transactions
		for _, t := range response {
			transaction := model.Transaction{
				TransactionID: t.TransactionId,
				CorporationID: ctx.corporationID,
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

			ctx.log.Debugf("Discovered new transaction %d (%d, %d x %.2f ISK))", transaction.TransactionID, transaction.TypeID, transaction.Quantity, transaction.UnitPrice)

			if err := db.InsertTransaction(&transaction); err != nil {
				ctx.log.Errorf("Could not insert transaction with ID %d: %v", transaction.TransactionID, err.Error())
			}
		}
	} else {
		ctx.log.WithFields(limitFields).Info("Transactions have not changed")
	}

	return httpResponse, nil
}
