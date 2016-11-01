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
	// c.AddFunc(gron.Every(1 * xtime.Day).At("10:00"),print)
	c.AddFunc(gron.Every(10*time.Second), birthday)
	c.Start()
}
