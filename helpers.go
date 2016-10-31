package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/roylee0704/gron"
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
func print() {
	_, monthnow, daynow := time.Now().Date()
	_, monthprev, dayprev := time.Now().Add(24 * time.Hour).Date()
	dates, err := getLines(*birthdayFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, date := range dates {
		if date.Month == monthnow.String() && daynow == date.Day {
			fmt.Printf("Wish birthday to %s", date.Name)
		} else if date.Month == monthprev.String() && dayprev == date.Day {
			fmt.Printf("Plan birthday for %s", date.Name)
		}
	}

}
func giveNotification() {
	c := gron.New()
	// c.AddFunc(gron.Every(1 * xtime.Day).At("10:00"),print)
	c.AddFunc(gron.Every(1*time.Second), print)
	c.Start()
}
