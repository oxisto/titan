/*
Copyright 2020 Christian Banse

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package titan

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/contracts"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/finance"
	"github.com/oxisto/titan/manufacturing"
	"github.com/oxisto/titan/model"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "app")
}

// App represented the main titan application
type App struct {
	CacheManufacturing bool
	CorporationID      int32
}

// ImportSDE reads the current SDE version from sde.version and imports it into the DB, if necessary.
func (a App) ImportSDE() {
	data, err := ioutil.ReadFile("sde.version")
	if err != nil {
		log.Error("Could not read SDE version, skipping import.")
		return
	}

	array := strings.Split(string(data), "-")
	if len(array) != 2 {
		log.Error("Could not read SDE version, skipping import.")
		return
	}

	i, err := strconv.Atoi(array[0])
	version := int32(i)
	server := array[1]
	if err != nil {
		log.Error("Could not read SDE version, skipping import.")
		return
	}

	log.Infof("Checking, if SDE %d is already cached...", version)

	sde := db.StaticDataExport{}
	cache.ReadCachedObject(fmt.Sprintf("sde:%d", version), &sde)

	if sde.Version == 0 {
		db.RunSDERestoreScript(version, array[1])
		sde.Version = version
		sde.Server = server

		cache.WriteCachedObject(sde)
	} else {
		log.Infof("SDE %d is already imported.", version)
	}
}

// ServerLoop takes care of reguarly caching prices and manufacturing.
func (a App) ServerLoop() {
	// builderID := int32(90821267)
	// builder := model.Character{}
	// cache.GetCharacter(builderID, &builder)
	typeIDs := []int32{}

	var productTypeIDs []int32
	var err error
	if productTypeIDs, err = db.GetProductTypeIDs(); err != nil {
		return
	}

	if !a.CacheManufacturing {
		return
	}

	typeIDs = append(typeIDs, productTypeIDs...)
	typeIDs = append(typeIDs, db.GetTech1BlueprintIDs()...)
	typeIDs = append(typeIDs, db.GetMaterialTypeIDs(manufacturing.ActivityManufacturing)...)
	typeIDs = append(typeIDs, db.GetMaterialTypeIDs(manufacturing.ActivityInvention)...)

	uniqueTypeIDs := MakeUnique(typeIDs)

	// this will cache all manufacturing objects, every hour
	for {
		log.Printf("Need to know the price of %d unique types.", len(uniqueTypeIDs))

		cache.GetPrices(model.JitaRegionID, uniqueTypeIDs)

		log.Printf("Trying to calculate profit for %d types...", len(productTypeIDs))

		for _, typeID := range productTypeIDs {
			/*go*/ a.UpdateProduct(typeID)
		}

		time.Sleep(time.Duration(1) * time.Hour)
	}
}

func (a App) ContractsLoop() {
	for {
		log.Printf("Trying get contracts for Jita region...")

		contracts.FetchContracts()

		time.Sleep(time.Duration(1) * time.Hour)
	}
}

func MakeUnique(slice []int32) []int32 {
	u := make([]int32, 0, len(slice))
	m := make(map[int32]bool)

	for _, v := range slice {
		if !m[v] {
			m[v] = true
			u = append(u, v)
		}
	}

	return u
}

func (a App) UpdateProduct(typeID int32) {
	m := model.Manufacturing{}

	if err := manufacturing.NewManufacturing(nil, int32(typeID), 10, 20, 0.1, &m); err == nil {
		db.UpdateProfit(m)
	} else {
		log.Printf("Error while manufacturing %s (%d): %v", m.Product.TypeName, typeID, err)
	}
}

func (a App) JournalLoop() {
	for {
		corporationID := a.CorporationID

		log.Printf("Fetching journal data for corporation %d...", corporationID)
		duration, err := finance.FetchJournal(corporationID, 1)

		if err != nil {
			log.Printf("An error occured while fetching journal: %v", err.Error())
		}

		log.Printf("Waiting for %.2f minutes until next fetch", duration.Minutes())

		time.Sleep(duration)
	}
}

func (a App) TransactionLoop() {
	for {
		corporationID := a.CorporationID

		log.Printf("Fetching transactions for corporation %d...", corporationID)
		duration, err := finance.FetchTransations(corporationID, 1)

		if err != nil {
			log.Printf("An error occured while fetching transactions: %v", err.Error())
		}

		log.Printf("Waiting for %.2f minutes until next fetch", duration.Minutes())

		time.Sleep(duration)
	}
}
