package main

import (
	"fmt"
	"time"
)

func generator(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Second)
		}
	}()

	return ch
}

func main() {
	ch := generator("Hello")
	for i := 0; i < 5; i++ {
		fmt.Println(<-ch)
	}
}
