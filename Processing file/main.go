package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	_ "os"
	"sort"
	"sync"
	"time"
)

type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

func processSequentially(people []Person) []Person {
	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
	return people
}

func processParallel(people []Person) []Person {
	var wg sync.WaitGroup
	result := make([]Person, len(people))

	ch := make(chan struct {
		index  int
		person Person
	}, len(people))

	for i := range people {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- struct {
				index  int
				person Person
			}{i, people[i]}
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for item := range ch {
		result[item.index] = item.person
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func readJSON(filename string) ([]Person, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var people []Person
	err = json.Unmarshal(file, &people)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func writeJSON(filename string, data []Person) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	people, err := readJSON("data.json")
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	sequentialResult := processSequentially(people)
	sequentialDuration := time.Since(start)

	err = writeJSON("sequential_result.json", sequentialResult)
	if err != nil {
		log.Fatal(err)
	}

	start = time.Now()
	parallelResult := processParallel(people)
	parallelDuration := time.Since(start)

	err = writeJSON("parallel_result.json", parallelResult)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sequential processing took: %v\n", sequentialDuration)
	fmt.Printf("Parallel processing took: %v\n", parallelDuration)
}
