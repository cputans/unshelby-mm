package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Bot struct {
	Token    string
	Server   string
	Team     *model.Team
	MMClient *model.Client4
	WSClient *model.WebSocketClient
}

func NewBot(token, server, teamName string) (*Bot, error) {
	bot := &Bot{
		Token:  token,
		Server: server,
	}

	bot.MMClient = model.NewAPIv4Client(fmt.Sprintf("https://%s", bot.Server))
	bot.MMClient.SetToken(bot.Token)

	team, _, _ := bot.MMClient.GetTeamByName(teamName, "")
	bot.Team = team

	return bot, nil
}

func (bot *Bot) Listen() {
	var err error
	wsUrl := fmt.Sprintf("wss://%s:443", bot.Server)
	log.Printf("connecting to %s", wsUrl)
	for {
		bot.WSClient, err = model.NewWebSocketClient4(
			wsUrl,
			bot.Token,
		)

		if err != nil {
			panic(err)
		}

		bot.WSClient.Listen()

		for event := range bot.WSClient.EventChannel {
			bot.HandleEvent(event)
		}
	}
}

func (bot *Bot) HandleEvent(event *model.WebSocketEvent) {
	if event.EventType() == model.WebsocketEventPosted {
		eventData := event.GetData()
		postJson := eventData["post"].(string)
		var post *model.Post
		json.Unmarshal([]byte(postJson), &post)
		if strings.Contains(post.Message, "setup") && eventData["sender_name"].(string) != "@unshelby" {
			rand := rand.IntN(6) + 1

			if rand == 3 {
				bot.MMClient.CreatePost(&model.Post{
					ChannelId: post.ChannelId,
					Message:   "Did somebody say setup?!?  https://denike.io/about/my-setup/desktop/",
				})
			}
		}
	}
}
