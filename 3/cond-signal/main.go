package main

import (
	"fmt"
	"sync"
)

const max = 10

func main() {
	c := sync.NewCond(&sync.Mutex{})

	queue := make([]string, 0, max)

	doSomething := func() {
		c.L.Lock()

		fmt.Println("doSomething")

		queue = queue[1:]
		fmt.Println("removed from queue")

		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < max+20; i++ {
		c.L.Lock()

		if len(queue) == max {
			fmt.Println("queue is full, waiting...")
			c.Wait()
		}

		queue = append(queue, fmt.Sprintf("element-%d", i))

		go doSomething()

		c.L.Unlock()
		c.Signal()
	}
}
