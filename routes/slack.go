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
