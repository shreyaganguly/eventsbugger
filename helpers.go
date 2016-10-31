package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/nlopes/slack"
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

func getLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		fmt.Println(scanner.Text())
		fmt.Println("**************")
	}
	return lines, scanner.Err()

}
