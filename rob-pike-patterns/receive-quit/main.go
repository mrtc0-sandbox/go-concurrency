package main

import "fmt"

func main() {
	quit := make(chan string)

	ch := generator("Hello", quit)

	for i := 0; i < 50; i++ {
		fmt.Println(<-ch, i)
	}

	quit <- "Done"
	// goroutine が終了するまでブロックする ≒ goroutine からの終了通知を待つ
	fmt.Println(<-quit)
}

func generator(msg string, quit chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case <-quit:
				quit <- "Bye Bye ~"
				return
			case ch <- fmt.Sprintf("%s", msg):
			}
		}
	}()

	return ch
}
