package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func worker(id int, word string, wg *sync.WaitGroup) {

	defer wg.Done()

	url := fmt.Sprintf(
		"http://localhost:8081/v1/anagrams?word=%s",
		word,
	)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("request error:", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	resp.Body.Close()

	fmt.Printf("worker %d -> %s\n", id, body)
}

func main() {

	words := []string{
		"read", "trace", "écart", "dear",
	}

	programArgs := os.Args[1:]
	var requests int = 20
	if len(programArgs) == 1 {
		i, err := strconv.Atoi(programArgs[0])

		if err != nil {
			fmt.Printf("First argument should be a number, default to 20")
		}
		requests = i
	}

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {

		wg.Add(1)

		go worker(i, words[i%len(words)], &wg)
	}

	wg.Wait()
}
