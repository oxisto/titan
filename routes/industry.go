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

func GetIndustryJobs(c *gin.Context) {
	character := c.Value(CharacterContext).(*model.Character)
	jobs := &model.IndustryJobs{}

	err := cache.GetIndustryJobs(character.CharacterID, character.CorporationID, jobs)

	JSON(c, http.StatusOK, jobs, err)
}
