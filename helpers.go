package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"

	"github.com/roylee0704/gron"
)

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

func giveNotification() {
	c := gron.New()
	// weekly := gron.Every(1 * xtime.Week)
	// c.AddFunc(gron.Every(1 * xtime.Day).At(*eventTime),print)
	c.AddFunc(gron.Every(10*time.Second), birthday)
	c.AddFunc(gron.Every(30*time.Second), github)
	// c.AddFunc(weekly, github)
	c.Start()
}
