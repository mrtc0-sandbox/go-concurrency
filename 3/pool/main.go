package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	myPool := &sync.Pool{
		New: func() interface{} {
			return "Hello"
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			myPool.Put("World")
			time.Sleep(300 * time.Millisecond)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 20; i++ {
			fmt.Println(myPool.Get())
			time.Sleep(100 * time.Millisecond)
		}
		wg.Done()
	}()

	wg.Wait()
}
