package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Repositories struct {
	Items []struct {
		FullName        string `json:"full_name"`
		HtmlURL         string `json:"html_url"`
		Description     string `json:"description"`
		StargazersCount int64  `json:"stargazers_count"`
	} `json:"items"`
}

func getRepositories() {
	v := url.Values{}
	v.Add("q", "language:go+created:>2016-09-26")
	v.Add("sort", "stars")
	v.Add("order", "desc")
	query, err := url.QueryUnescape(v.Encode())
	if err != nil {
		log.Fatal("ERROR******", err)
	}
	url := fmt.Sprintf("https://api.github.com/search/repositories?%s", query)
	resp, err := http.Get(url)
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ERROR******", err)
	}
	var r Repositories
	err = json.Unmarshal(htmlData, &r)

	if err != nil {
		log.Fatal("ERROR********", err)
	}
	defer resp.Body.Close()
	fmt.Println("REsponse******************", r)

}
