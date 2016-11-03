package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type BirthDay struct {
	Name  string
	Month string
	Day   int
}

func getBirthdays(filepath string) ([]BirthDay, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var birthDays []BirthDay
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.Split(scanner.Text(), ":")
		splits := strings.Split(scanner.Text(), " ")
		day, err := strconv.Atoi(splits[1])
		if err != nil {
			return nil, err
		}
		birthDays = append(birthDays, BirthDay{name[0], splits[2], day})
	}
	return birthDays, scanner.Err()

}

func postBirthdayMessage(members []slack.IM, bID, message string) {
	for _, member := range members {
		if member.User != bID {
			// chanID, timestamp, err := client.PostMessage(member.ID, message, slack.PostMessageParameters{})
			// if err != nil {
			// 	fmt.Printf("%s\n", err)
			// 	return
			// }
			fmt.Printf("Message successfully sent to channel %s at %s", channelID, "123")
		} else {
			params := slack.PostMessageParameters{
				UnfurlLinks: true,
				UnfurlMedia: true,
				Text:        "Regards, Team Paypermint",
			}
			chanID, timestamp, err := client.PostMessage(member.ID, getRandomBirthdayURL()+"  Regards, Team Paypermint", params)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("Message successfully sent to channel %s at %s", chanID, timestamp)
		}
	}
}

func birthday() {
	_, monthnow, daynow := time.Now().Date()
	_, monthprev, dayprev := time.Now().Add(24 * time.Hour).Date()
	dates, err := getBirthdays(*birthdayFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, date := range dates {
		if date.Month == monthnow.String() && daynow == date.Day {
			message := fmt.Sprintf("Wish birthday to %s", date.Name)
			birthdayID, err := getUserID(date.Name)
			if err != nil {
				if err.Error() == "No userID found" {
					message = fmt.Sprintf("%s is not on slack, give him/her a call", date.Name)
					chanID, timestamp, errpost := client.PostMessage(channelID, message, slack.PostMessageParameters{})
					if errpost != nil {
						fmt.Printf("%s\n", errpost)
						return
					}
					fmt.Printf("Message successfully sent to channel %s at %s", chanID, timestamp)
					return
				} else {
					fmt.Printf("%s\n", err)
					return
				}

			}
			members, err := getMembers()
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			postBirthdayMessage(members, birthdayID, message)
		} else if date.Month == monthprev.String() && dayprev == date.Day {
			message := fmt.Sprintf("Plan birthday for %s", date.Name)
			members, err := getMembers()
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			birthdayID, err := getUserID(date.Name)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			postBirthdayMessage(members, birthdayID, message)
		}
	}

}
