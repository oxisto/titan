// Package datafetch contains code that fetch various data, such as transactions, jour«nals, etc. of a company
package datafetch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type DataFetcher interface {
	// Fetch is the main entrypoint for the data fetcher. It should contain code that fetches data from
	// an external source, such as ESI. It needs to return the HTTP response, so that the fetch service
	// can check the response for caching headers, such as ETag or expire. Ideally, a fetcher should only
	// contain one HTTP request. If more than one request is issued, the one with the longest cache time
	// should be returned.
	//
	// ctx contains a FetchContext, which holds all necessary information the fetcher needs about the
	// current fetch request, such as the corporation ID and a logger.
	Fetch(ctx FetchContext) (*http.Response, error)

	// DataType specifies the data type that is returned by this fetcher.
	DataType() string

	// MaxCacheTime specifies the cache time we should uphold, if we cannot get a better (possibly smaller)
	// value. This should be conservative to not bother the external service too much.
	MaxCacheTime() time.Duration

	// LogFields can be used to inject additional fields into the logger supplied by the FetchContext.
	// The primary use case is to display additional values, that are specific to a fetcher. For an example,
	// have a look at the transactionFetcher, that has a specific division field.
	LogFields() logrus.Fields
}

type FetchService struct {
	corporationID int32
	fetcher       DataFetcher
}

type FetchContext struct {
	context.Context

	corporationID int32
	lastETag      string
	accessToken   *model.AccessToken
	log           *logrus.Entry
}

func init() {
	log = logrus.WithField("component", "data-fetcher")
}

type CachedETag struct {
	ETag string
	Key  string
}

func (t CachedETag) ID() int32 {
	return 0
}

func (t CachedETag) ExpiresOn() *time.Time {
	return nil
}

func (t *CachedETag) SetExpire(time *time.Time) {

}

func (t CachedETag) HashKey() string {
	return fmt.Sprintf("etag:%s", t.Key)
}

type metadata struct {
	dataType         string
	maxCacheTime     time.Duration
	additionalFields logrus.Fields
}

func (m metadata) DataType() string {
	return m.dataType
}

func (m metadata) MaxCacheTime() time.Duration {
	return m.maxCacheTime
}

func (m metadata) LogFields() logrus.Fields {
	return m.additionalFields
}

func NewFetchService(corporationID int32, fetcher DataFetcher) *FetchService {
	return &FetchService{
		fetcher:       fetcher,
		corporationID: corporationID,
	}
}

func (service FetchService) StartLoop() {
	var (
		backoffTime = time.Minute
		etag        CachedETag
		accessToken model.AccessToken
		ctx         FetchContext
		err         error
	)

	// create a new context
	ctx = FetchContext{
		corporationID: service.corporationID,
		log: log.WithFields(logrus.Fields{
			"data":          service.fetcher.DataType(),
			"corporationID": service.corporationID,
		}).WithFields(service.fetcher.LogFields()),
	}

	for {
		if ctx.lastETag == "" {
			// check, if we have an ETag in our cache
			_ = cache.ReadCachedObject(fmt.Sprintf("etag:%s", cacheKey(ctx, service.fetcher)), &etag)
			if err != nil {
				// just warn and ignore cache errors, because they do not influence our fetching
				log.Warnf("Could not read ETag: %v", err)
			} else {
				ctx.lastETag = etag.ETag
			}
		}

		ctx.log.Printf("Fetching %s...", service.fetcher.DataType())

		// find access token for corporation
		err = cache.GetAccessTokenForCorporation(ctx.corporationID, &accessToken)
		if err != nil {
			ctx.log.Errorf("Could not find access token for %d.", ctx.corporationID, backoffTime.Minutes())

			// this error could occur, if no access tokens are ready yet. let's wait for a little bit
			sleepWithPrint(ctx.log, backoffTime)

			// increase the backoff time
			backoffTime = backoffTime * 2

			continue
		}

		// update the context with the access token
		ctx.accessToken = &accessToken

		// let the fetcher do its work. it will return the time we need to sleep
		httpResponse, err := service.fetcher.Fetch(ctx)
		if err != nil {
			ctx.log.Printf("An error occured while fetching %s: %v", service.fetcher.DataType(), err)

			// just to be sure, we wait for the maximum time to not bother ESI too much
			sleepWithPrint(ctx.log, service.fetcher.MaxCacheTime())
			continue
		}

		// try to parse the expires header from the http response
		expireTime, err := time.Parse(time.RFC1123, httpResponse.Header.Get("expires"))
		if err != nil {
			ctx.log.Printf("An error occured while parsing the expires header: %v", err)

			// just to be sure, we wait for the maximum time to not bother ESI too much
			sleepWithPrint(ctx.log, service.fetcher.MaxCacheTime())
			continue
		}

		etag = CachedETag{
			ETag: httpResponse.Header.Get("etag"),
			Key:  cacheKey(ctx, service.fetcher),
		}

		// update the ETag directly
		ctx.lastETag = etag.ETag

		// cache the ETag
		err = cache.WriteCachedObject(&etag)
		if err != nil {
			// just warn and ignore cache errors, because they do not influence our fetching
			log.Warnf("Could not cache ETag: %v", err)
		}

		var duration = time.Until(expireTime)

		// sometimes, the duration is negative, this can occur because of clock offset between the server and our client.
		// in this case we need to assume the maximum time
		if duration < 0 {
			duration = service.fetcher.MaxCacheTime()
		}

		sleepWithPrint(ctx.log, duration)
	}
}

func sleepWithPrint(log *logrus.Entry, duration time.Duration) {
	log.Printf("Waiting for %.2f minutes until next fetch", duration.Minutes())
	time.Sleep(duration)
}

func cacheKey(ctx FetchContext, fetcher DataFetcher) string {
	return fmt.Sprintf("%s:%d:%v", fetcher.DataType(), ctx.corporationID, fetcher.LogFields())
}
