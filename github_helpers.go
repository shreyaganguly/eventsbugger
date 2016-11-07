package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type Repositories struct {
	Items []struct {
		Name            string `json:"name"`
		HtmlURL         string `json:"html_url"`
		StargazersCount int64  `json:"stargazers_count"`
	} `json:"items"`
}

func github() {
	v := url.Values{}
	d := strings.Split(time.Now().Add(-24*7*time.Hour).String(), " ")[0]
	queryString := fmt.Sprintf("language:go+created:>%s", d)
	v.Add("q", queryString)
	v.Add("sort", "stars")
	v.Add("order", "desc")
	query, err := url.QueryUnescape(v.Encode())
	if err != nil {
		log.Fatal("ERROR!!", err)
	}
	url := fmt.Sprintf("https://api.github.com/search/repositories?%s", query)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ERROR!!", err)
	}
	var r Repositories
	err = json.Unmarshal(htmlData, &r)

	if err != nil {
		log.Fatal("ERROR!!", err)
	}
	var message string
	if len(r.Items) > 0 {
		chanID, timestamp, errpost := client.PostMessage(gchannelID, "Your Weekly Dose Of Go Fever", slack.PostMessageParameters{})
		if errpost != nil {
			fmt.Printf("%s\n", errpost)
			return
		}
		fmt.Printf("Message successfully sent to channel %s at %s", chanID, timestamp)
	}
	for i, repository := range r.Items {
		if i == 5 {
			break
		}
		// <http://www.zapier.com|Text to make into a link>
		message = fmt.Sprintf("<%s|%s> is trending with %d star gazers\n", repository.HtmlURL, strings.ToUpper(repository.Name), repository.StargazersCount)
		// message = fmt.Sprintf("%s\n%s", message, m)
		params := slack.PostMessageParameters{
			Text:        repository.HtmlURL,
			UnfurlLinks: true,
			UnfurlMedia: true,
		}
		chanID, timestamp, errpost := client.PostMessage(gchannelID, message, params)
		if errpost != nil {
			fmt.Printf("%s\n", errpost)
			return
		}
		fmt.Printf("Message successfully sent to channel %s at %s", chanID, timestamp)
	}
	// params := slack.PostMessageParameters{
	// 	UnfurlLinks: true,
	// 	UnfurlMedia: true,
	// }
	// chanID, timestamp, errpost := client.PostMessage(gchannelID, message, params)
	// if errpost != nil {
	// 	fmt.Printf("%s\n", errpost)
	// 	return
	// }
	// fmt.Printf("Message successfully sent to channel %s at %s", chanID, timestamp)
	fmt.Println("response", r)

}
