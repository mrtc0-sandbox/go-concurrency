package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web    = fakeSearch("web")
	Web2   = fakeSearch("web2")
	Image  = fakeSearch("image")
	Image2 = fakeSearch("image2")
	Video  = fakeSearch("video")
	Video2 = fakeSearch("video2")
)

type Result string
type Search func(query string) Result

func main() {
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}

	return <-c

}

func Google(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- First(query, Web, Web2) }()
	go func() { c <- First(query, Image, Image2) }()
	go func() { c <- First(query, Video, Video2) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}

	return
}
