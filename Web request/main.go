package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Stats struct {
	UserCount  int
	ErrorCount int
}

var apiEndpoints = []string{
	"https://jsonplaceholder.typicode.com/users",
	"https://jsonplaceholder.typicode.com/posts",
	"https://jsonplaceholder.typicode.com/comments",
}

func getData(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: Status code %d", resp.StatusCode)
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Error decoding response: %v", err)
	}

	return data, nil
}

func getStatsSequential() Stats {
	var stats Stats

	for _, url := range apiEndpoints {
		data, err := getData(url)
		if err != nil {
			stats.ErrorCount++
			fmt.Printf("Error from API: %s\n", err)
		} else {
			switch v := data.(type) {
			case []interface{}:
				stats.UserCount += len(v)
				fmt.Printf("Data fetched: %d records\n", len(v))
			}
		}
	}

	return stats
}

func getStatsParallel() Stats {
	var wg sync.WaitGroup
	ch := make(chan interface{}, len(apiEndpoints))
	var stats Stats

	for _, url := range apiEndpoints {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			data, err := getData(url)
			if err != nil {
				ch <- fmt.Sprintf("Error from API: %s", err)
			} else {
				ch <- data
			}
		}(url)
	}

	wg.Wait()
	close(ch)

	for data := range ch {
		switch v := data.(type) {
		case string:
			stats.ErrorCount++
			fmt.Printf("%s\n", v)
		case []interface{}:
			stats.UserCount += len(v)
			fmt.Printf("Data fetched: %d records\n", len(v))
		}
	}

	return stats
}

func main() {
	start := time.Now()
	statsSeq := getStatsSequential()
	elapsedSeq := time.Since(start)

	start = time.Now()
	statsPar := getStatsParallel()
	elapsedPar := time.Since(start)

	fmt.Printf("\nSequential Execution:\n")
	fmt.Printf("Total users fetched: %d\n", statsSeq.UserCount)
	fmt.Printf("Number of errors: %d\n", statsSeq.ErrorCount)
	fmt.Printf("Time taken: %v\n", elapsedSeq)

	fmt.Printf("\nParallel Execution:\n")
	fmt.Printf("Total users fetched: %d\n", statsPar.UserCount)
	fmt.Printf("Number of errors: %d\n", statsPar.ErrorCount)
	fmt.Printf("Time taken: %v\n", elapsedPar)
}
