package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/reconquest/karma-go"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("TOKEN is not specified")
	}

	channel := os.Getenv("CHANNEL")
	if channel == "" {
		log.Fatalln("CHANNEL is not specified")
	}

	word := os.Getenv("WORD")
	if word == "" {
		log.Fatalln("WORD is not specified")
	}

	word = strings.ToLower(word)

	api := slack.New(
		token,
		slack.OptionLog(
			log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags),
		),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.Println("Successfully connected to Slack WS")

		case *slack.MessageEvent:
			err := handleMessage(api, event, channel, word)
			if err != nil {
				log.Fatalln(err)
			}

		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials")
			return
		}
	}
}

func handleMessage(
	api *slack.Client,
	event *slack.MessageEvent,
	channel string,
	word string,
) error {
	if event.Channel != channel {
		return nil
	}

	if strings.Contains(strings.ToLower(event.Text), word) {
		return sendChart(api, event.Channel)
	}

	return nil
}

func sendChart(api *slack.Client, channel string) error {
	text := "http://j1.profinance.ru/delta/prochart?type=USDRUB&" +
		"amount=60&chart_height=220&chart_width=400&" +
		"grtype=2&tictype=3&iId=5&seed=" + fmt.Sprint(time.Now().UnixNano())

	_, _, err := api.PostMessage(channel, slack.MsgOptionText(text, false))
	if err != nil {
		return karma.Format(
			err,
			"unable to post message",
		)
	}

	return nil
}
