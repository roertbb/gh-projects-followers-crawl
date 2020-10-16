package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const ttl = 3
const skipAbove = 30

// PeopleWithTTL ...
type PeopleWithTTL struct {
	People []Person
	TTL    int
}

func main() {
	username := os.Args[1]
	searchPhrase := os.Args[2]

	// person/ttl
	visitedPeople := make(map[Person]int)

	// htmlUrl/repo
	repos := make(map[string]Repo)

	peopleChan := make(chan PeopleWithTTL)
	reposChan := make(chan []Repo)

	// TODO: proper signaling when all worker finished - pass message via channel
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)

		for {
			select {
			case data := <-peopleChan:
				for _, person := range data.People {
					if _, present := visitedPeople[person]; !present && data.TTL > 0 {
						visitedPeople[person] = data.TTL
						go getReposForUser(person.Login, reposChan)
						go getFollowersForUser(person.Login, peopleChan, data.TTL-1)
						go getFollowingForUser(person.Login, peopleChan, data.TTL-1)
					}
				}
			case data := <-reposChan:
				for _, repo := range data {
					if strings.Contains(repo.Name, searchPhrase) || strings.Contains(repo.Description, searchPhrase) {
						repos[repo.HtmlUrl] = repo
						fmt.Println(repo)
					}
				}
			}
		}

	}()

	peopleChan <- PeopleWithTTL{People: []Person{{Login: username}}, TTL: ttl}

	wg.Wait()
}
