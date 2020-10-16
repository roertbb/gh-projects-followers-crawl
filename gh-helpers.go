package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func basicAuth() string {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth())
	return nil
}

// Repo ...
type Repo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
}

func getReposForUser(username string, reposChan chan<- []Repo) {
	url := "https://api.github.com/users/" + username + "/repos"

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth())
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	repos := make([]Repo, 0)
	jsonErr := json.Unmarshal([]byte(body), &repos)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	reposChan <- repos
}

// Person ...
type Person struct {
	Login string `json:"login"`
}

func getFollowersForUser(username string, peopleChan chan<- PeopleWithTTL, ttl int) {
	url := "https://api.github.com/users/" + username + "/followers"
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth())
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	followers := make([]Person, 0)
	jsonErr := json.Unmarshal([]byte(body), &followers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	peopleChan <- PeopleWithTTL{People: followers, TTL: ttl}
}

func getFollowingForUser(username string, peopleChan chan<- PeopleWithTTL, ttl int) {
	url := "https://api.github.com/users/" + username + "/following"
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth())
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	following := make([]Person, 0)
	jsonErr := json.Unmarshal([]byte(body), &following)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	peopleChan <- PeopleWithTTL{People: following, TTL: ttl}
}
