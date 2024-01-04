package main

import (
	"fmt"
)

type Message struct {
	str   string
	block chan int
}

// 多重化しつつ、順序を保つ例
func main() {
	ch := fanIn(generator("Hello"), generator("World"))

	for i := 0; i < 10; i++ {
		msg1 := <-ch
		fmt.Println(msg1.str)

		msg2 := <-ch
		fmt.Println(msg2.str)

		// ここでブロックしているので、順序を保つことができる
		<-msg1.block
		<-msg2.block
	}
}

func fanIn(ch1, ch2 <-chan Message) <-chan Message {
	ch := make(chan Message)

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

func generator(msg string) <-chan Message {
	ch := make(chan Message)
	blockingStep := make(chan int)

	go func() {
		for i := 0; ; i++ {
			ch <- Message{fmt.Sprintf("%s %d", msg, i), blockingStep}

			blockingStep <- 1
		}
	}()

	return ch
}
