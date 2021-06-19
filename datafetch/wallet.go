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

type walletFetcher struct {
	corporationID int32

	log *logrus.Entry
}

func NewWallerFetcher(corporationID int32, division int32) DataFetcher {
	return &walletFetcher{
		corporationID: corporationID,
		log: log.WithFields(logrus.Fields{
			"data":          "wallet",
			"corporationID": corporationID,
		}),
	}
}

func (f walletFetcher) StartLoop() {
	for {
		f.log.Printf("Fetching wallet...")
		duration, err := f.Fetch()

		if err != nil {
			f.log.Errorf("An error occured while fetching wallet: %v", err)
		}

		if duration < 0 {
			duration = f.MaxCacheTime()
		}

		f.log.Printf("Waiting for %.2f minutes until next fetch", duration.Minutes())

		time.Sleep(duration)
	}
}

func (f walletFetcher) MaxCacheTime() time.Duration {
	return time.Minute * 5
}

func (f walletFetcher) Fetch() (time.Duration, error) {
	// find access token for corporation
	accessToken := model.AccessToken{}
	err := cache.GetAccessTokenForCorporation(f.corporationID, &accessToken)
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	var options esi.GetCorporationsCorporationIdWalletsOpts

	response, httpResponse, err := cache.ESI.WalletApi.GetCorporationsCorporationIdWallets(
		context.WithValue(context.Background(),
			goesi.ContextAccessToken,
			accessToken.Token),
		f.corporationID,
		&options)
	if err != nil {
		return time.Minute * 5, err
	}

	t, err := time.Parse(time.RFC1123, httpResponse.Header.Get("Expires"))
	if err != nil {
		return time.Minute * 5, err
	}

	// loop through all transactions
	for _, t := range response {
		wallet := model.Wallet{
			CorporationID: int64(f.corporationID),
			Division:      t.Division,
			Balance:       t.Balance,
		}

		f.log.Printf("Updating balance for division %d (%.2f ISK)", wallet.Division, wallet.Balance)

		if err := db.UpdateWallet(&wallet); err != nil {
			f.log.Errorf("Could not update wallet: %v", err)
		}
	}

	return time.Until(t), nil
}
