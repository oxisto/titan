/*
Copyright 2018 Christian Banse

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

package slack

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/oxisto/titan/cache"
	"github.com/oxisto/titan/manufacturing"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Something struct {
	Channel string
	Msg     slack.Msg
}

var (
	replyChannel chan Something
	api          *slack.Client
	botId        string
	log          *logrus.Entry
)

func init() {
	log = logrus.WithField("component", "slack")
}

func Bot(token string) {
	if token != "" {
		log.Info("Connecting to Slack...")

		api = slack.New(token)

		replyChannel = make(chan Something)
		go handleBotReply()

		//api.SendMessage("#eve", slack.MsgOptionText("Titan Server has started.", false), slack.MsgOptionPost(), slack.MsgOptionAsUser(true))

		rtm := api.NewRTM()
		go rtm.ManageConnection()

		for {
			select {
			case msg := <-rtm.IncomingEvents:
				switch ev := msg.Data.(type) {
				case *slack.ConnectedEvent:
					botId = ev.Info.User.ID
					log.Infof("Connected to Slack using %s", botId)
				case *slack.TeamJoinEvent:
					// Handle new user to client
				case *slack.MessageEvent:
					something := Something{
						Msg:     ev.Msg,
						Channel: ev.Channel,
					}

					if ev.Msg.User != botId {
						replyChannel <- something
					}
				case *slack.ReactionAddedEvent:
					// Handle reaction added
				case *slack.ReactionRemovedEvent:
					// Handle reaction removed
				case *slack.RTMError:
					log.Errorf("Error: %s", ev.Error())
				}
			}
		}
	}
}

func handleBotReply() {
	for {
		something := <-replyChannel

		text := something.Msg.Text

		if strings.Contains(text, "profit") {
			//api.SendMessage(something.Channel, slack.MsgOptionText("You are interested in profit?", false), slack.MsgOptionPost(), slack.MsgOptionAsUser(true))
			log.Debug("Triggering profit command...")
			profitCommand(something)
		}
	}
}

func profitCommand(something Something) {
	marks := []string{",", ".", "!", ":", "?"}
	fillerWords := map[string]bool{"of": true, "a": true, "the": true, "please": true}

	command := something.Msg.Text

	// first, remove all punctuation marks
	for _, mark := range marks {
		command = strings.Replace(command, mark, "", -1)
	}

	// next, tokenize it and remove all filler words and re-assemble the rest after the keyword was found
	tokens := strings.Split(command, " ")
	typeNameTokens := []string{}
	keywordFound := false
	keyword := "profit"
	for _, token := range tokens {
		if token == keyword {
			keywordFound = true
			continue
		}

		if fillerWords[token] {
			continue
		}

		if keywordFound {
			typeNameTokens = append(typeNameTokens, token)
		}
	}

	typeName := strings.Join(typeNameTokens, " ")

	builderID := int32(92925923)
	builder := model.Character{}
	cache.GetCharacter(builderID, &builder)

	// try to find it using the get product types
	options := cache.NewSearchOptions()
	options.NameFilter = typeName
	options.Limit = 1
	types, err := cache.GetProductTypes(options, builder)
	if err != nil {
		replyWithError(something.Channel, err)
		return
	}

	if len(types) == 0 {
		replyWithError(something.Channel, errors.New("item not found"))
		return
	}

	// pick the first type
	typeID := types[0].TypeID

	m := manufacturing.Manufacturing{}
	manufacturing.NewManufacturing(builder, int32(typeID), &m)

	p := message.NewPrinter(language.English)

	fields := []slack.AttachmentField{
		{
			Value: fmt.Sprintf("The %s is a %s manufactured from a %s.",
				m.Product.Name.EN,
				m.Product.Group.Name.EN,
				m.BlueprintType.Name.EN),
		},
		{
			Title: "Needs invention",
			Value: p.Sprintf("%t", m.IsTech2),
		},
		{
			Title: "Daily profit (based on sell orders)",
			Value: p.Sprintf("%.2f ISK.", m.Profit.PerDay.BasedOnSellPrice),
			Short: true,
		}, {
			Title: "Daily profit (based on buy orders)",
			Value: p.Sprintf("%.2f ISK.", m.Profit.PerDay.BasedOnBuyPrice),
			Short: true,
		}}

	/*actions := []slack.AttachmentAction{
		{
			Name: "open",
			Text: "Open in Market Browser",
			Type: "button",
		},
	}*/

	attachment := slack.Attachment{
		Color:     "#B733FF",
		Title:     fmt.Sprintf("%s", m.Product.Name.EN),
		TitleLink: fmt.Sprintf("https://eve.aybaze.com/#/manufacturing/%d", m.Product.TypeID),
		Fields:    fields,
		//Actions:    actions,
		ThumbURL:   fmt.Sprintf("https://image.eveonline.com/Type/%d_64.png", m.Product.TypeID),
		CallbackID: fmt.Sprintf("profit:%d", m.Product.TypeID),
	}

	params := slack.PostMessageParameters{}
	params.AsUser = true
	params.Attachments = []slack.Attachment{attachment}

	api.PostMessage(something.Channel, "", params)

	//api.SendMessage(something.Channel, slack.MsgOptionText(fmt.Sprintf("%s has a  See https://eve.aybaze.com/#/manufacturing/%d", m.Product.Name.EN, m.Profit.PerDay.BasedOnBuyPrice, m.Product.TypeID), false), slack.MsgOptionPost(), slack.MsgOptionAsUser(true))
}

func replyWithError(channel string, err error) {
	api.SendMessage(channel, slack.MsgOptionText(fmt.Sprintf("Does not compute: %s", err), false), slack.MsgOptionPost(), slack.MsgOptionAsUser(true))
}
