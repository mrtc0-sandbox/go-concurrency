package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	c := sync.NewCond(&sync.Mutex{})

	subscribe := func(i int) {
		c.L.Lock()
		defer c.L.Unlock()

		fmt.Println("waiting for broadcast in goroutine, i=", i)
		c.Wait()
		fmt.Println("received signal from main goroutine, i=", i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			subscribe(i)
		}(i)
	}

	time.Sleep(1 * time.Second)
	c.Broadcast()
	wg.Wait()
}
