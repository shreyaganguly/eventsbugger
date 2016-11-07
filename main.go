package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

var (
	slackBotToken       = flag.String("slackbottoken", "", "Token for the slack bot that is to be integrated with")
	birthdaychannelName = flag.String("bchannelname", "", "channel name where slack bot is to be integrated")
	githubchannelName   = flag.String("gchannelname", "", "channel name where slack bot is to be integrated for getting github updates")
	birthdayFile        = flag.String("birthdayfile", "", "the file that contains birthdaypersonname: Day Month(birthday) in each line")
	eventTime           = flag.String("eventTime", "", "Set a time of when to remind you of your events")
)

func main() {
	flag.Parse()
	api := slack.New(*slackBotToken)
	// logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	// slack.SetLogger(logger)
	bchanID, err := getChannelID(api, *birthdaychannelName)
	if err != nil {
		log.Fatal(err)
	}
	gchanID, err := getChannelID(api, *githubchannelName)
	if err != nil {
		log.Fatal(err)
	}
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				a := msg.Data.(*slack.MessageEvent)
				fmt.Printf("Message: %s\n", a.Msg.Text)
				fmt.Printf("Message: %v\n", ev)
				if a.Msg.Text == "birthday" {
					setSlackClient(api, bchanID, gchanID)
					giveNotification()
				}
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
