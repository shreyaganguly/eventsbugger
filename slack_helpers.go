package main

import (
	"errors"

	"github.com/nlopes/slack"
)

var (
	client    *slack.Client
	channelID string
)

func setSlackClient(c *slack.Client, cID string) {
	client = c
	channelID = cID
}

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
