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
)

func SlackCallback(w http.ResponseWriter, r *http.Request) {
	payload := r.Form.Get("payload")

	log.Debugf("Payload: %s\n", payload)

	/*_ := slack.AttachmentActionCallback{}

	err := json.Unmarshal([]byte(payload), &callback)

	// TODO: verify token

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	var idx int
	if idx = strings.Index(callback.CallbackID, "profit:"); idx == -1 {
		fmt.Printf("Unknown callback from slack, ignoring")
		return
	}

	characterID := int32(92925923)

	typeID, err := strconv.Atoi(callback.CallbackID[idx:])
	if err != nil {
		fmt.Printf("Could not convert %s into a typeID\n", callback.CallbackID[idx:])
	}

	OpenMarket(characterID, int32(typeID), w, r)

	defer r.Body.Close()*/

}
