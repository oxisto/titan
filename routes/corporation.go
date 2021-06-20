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

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/model"
)

func GetCorporation(c *gin.Context) {
	character := c.Value(CharacterContext).(*model.Character)
	corporation := &model.Corporation{}

	err := cache.GetCorporation(character.CharacterID, character.CorporationID, corporation)

	JSON(c, http.StatusOK, corporation, err)
}

func GetCorporationWallets(c *gin.Context) {
	character := c.Value(CharacterContext).(*model.Character)
	wallets := &model.Wallets{}

	err := cache.GetCorporationWallets(character.CharacterID, character.CorporationID, wallets)

	JSON(c, http.StatusOK, wallets, err)
}
