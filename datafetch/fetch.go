// Package datafetch contains code that fetch various data, such as transactions, jour«nals, etc. of a company
package datafetch

import (
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type DataFetcher interface {
	StartLoop()
	Fetch() (time.Duration, error)

	MaxCacheTime() time.Duration
}

func init() {
	log = logrus.WithField("component", "data-fetcher")
}
