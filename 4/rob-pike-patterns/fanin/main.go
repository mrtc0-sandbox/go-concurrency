package main

import (
	"fmt"
)

func generator(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
		}
	}()

	return ch
}

// fanIn は、複数のチャネルが一つのチャネルを通るようにする
// それぞれのチャネルからリダイレクトすることで multiplexing している
func fanIn(ch1, ch2 <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			ch <- <-ch1
		}
	}()
	go func() {
		for {
			ch <- <-ch2
		}
	}()

	return ch
}

func main() {
	ch := fanIn(generator("Hello"), generator("World"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}
