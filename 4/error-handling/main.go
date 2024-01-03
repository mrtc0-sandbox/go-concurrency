package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Err      error
	Response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{Err: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://example.com", "https://badhost", "https://blog.ssrf.in"}
	for result := range checkStatus(done, urls...) {
		if result.Err != nil {
			fmt.Printf("error: %v\n", result.Err)
			continue
		}

		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
