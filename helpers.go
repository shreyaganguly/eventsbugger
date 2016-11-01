package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/roylee0704/gron"
)

var (
	client    *slack.Client
	channelID string
)

func getChannelID(client *slack.Client, channelname string) (string, error) {
	channels, err := client.GetChannels(false)
	if err != nil {
		return "", err
	}
	for _, v := range channels {
		if v.Name == channelname {
			return v.ID, nil
		}
	}
	return "", errors.New("No channel id with the channel name found")
}

type BirthDay struct {
	Name  string
	Month string
	Day   int
}

func getLines(filepath string) ([]BirthDay, error) {
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
func setSlackClient(c *slack.Client, cID string) {
	client = c
	channelID = cID
}
func getMembers() ([]slack.IM, error) {
	IMChannels, err := client.GetIMChannels()
	if err != nil {
		return nil, err
	}
	return IMChannels, nil
}
func getUserID(name string) (string, error) {
	users, err := client.GetUsers()
	if err != nil {
		return "", err
	}
	for _, user := range users {
		if user.Name == name {
			return user.ID, nil
		}
	}
	return "", errors.New("No userID found")
}
func getRandomBirthdayURL() string {
	file, _ := os.Open("birthday/birthdaylinks.txt")
	defer file.Close()
	var birthdayURLs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		birthdayURLs = append(birthdayURLs, scanner.Text())
	}
	return birthdayURLs[rand.Intn(len(birthdayURLs))]
}
func postMessage(members []slack.IM, bID, message string) {
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
			}
			chanID, timestamp, err := client.PostMessage(member.ID, getRandomBirthdayURL(), params)
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
	dates, err := getLines(*birthdayFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, date := range dates {
		if date.Month == monthnow.String() && daynow == date.Day {
			message := fmt.Sprintf("Wish birthday to %s", date.Name)
			birthdayID, err := getUserID(date.Name)
			if err.Error() == "No userID found" {
				message = fmt.Sprintf("%s is not on slack, give him/her a call", date.Name)
				chanID, timestamp, errpost := client.PostMessage(channelID, message, slack.PostMessageParameters{})
				if errpost != nil {
					fmt.Printf("%s\n", errpost)
					return
				}
				fmt.Printf("Message successfully sent to channel %s at %s", chanID, timestamp)
				return
			}
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			members, err := getMembers()
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			postMessage(members, birthdayID, message)
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
			postMessage(members, birthdayID, message)
		}
	}

}
func giveNotification() {
	c := gron.New()
	// c.AddFunc(gron.Every(1 * xtime.Day).At("10:00"),print)
	c.AddFunc(gron.Every(10*time.Second), birthday)
	c.Start()
}
