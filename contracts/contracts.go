package contracts

import (
	"context"
	"log"
	"time"

	"titan/cache"
	"titan/model"

	"github.com/antihax/goesi/esi"
)

func FetchContracts() (time.Duration, error) {
	options := esi.GetContractsPublicRegionIdOpts{}

	response, httpResponse, err := cache.ESI.ContractsApi.GetContractsPublicRegionId(context.Background(), model.JitaRegionID, &options)

	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	t, err := time.Parse(time.RFC1123, httpResponse.Header.Get("Expires"))
	if err != nil {
		return time.Duration(3600 * 1000 * 1000), err
	}

	for _, contract := range response {
		log.Printf("%v", contract)
	}

	return t.Sub(time.Now()), nil
}
