package main

import "fmt"

func main() {
	quit := make(chan bool)

	ch := generator("Hello", quit)

	for i := 0; i < 50; i++ {
		fmt.Println(<-ch, i)
	}

	// goroutine を終了させる。終了の完了は待たない
	quit <- true
}

func generator(msg string, quit chan bool) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case <-quit:
				fmt.Println("Done")
				return
			case ch <- fmt.Sprintf("%s", msg):
			}
		}
	}()

	return ch
}
