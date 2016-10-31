package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nlopes/slack"
)

var (
	slackBotToken = flag.String("slackbottoken", "", "Token for the slack bot that is to be integrated with")
	channelName   = flag.String("channelname", "", "channel name where slack bot is to be integrated")
	birthdayFile  = flag.String("birthdayfile", "", "the file that contains slackusername: dd/mm/yyyy(birthday) in each line")
)

func getPugURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting response: ", err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Read response error: ", err)
	}
	puggy := make(map[string]string)
	err = json.Unmarshal(bytes, &puggy)
	if err != nil {
		fmt.Println("Json unmarshal error: ", err)
	}
	return puggy["pug"]
}
func main() {
	flag.Parse()
	api := slack.New(*slackBotToken)
	// logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	// slack.SetLogger(logger)
	chanID, err := getChannelID(api, *channelName)
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
				if a.Msg.Text == "pug me" {
					params := slack.PostMessageParameters{
						UnfurlLinks: true,
						UnfurlMedia: true,
					}
					channelID, timestamp, err := api.PostMessage(chanID, getPugURL("http://pugme.herokuapp.com/random"), params)
					if err != nil {
						fmt.Printf("%s\n", err)
						return
					}
					fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
				}
				if a.Msg.Text == "birthday" {
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
