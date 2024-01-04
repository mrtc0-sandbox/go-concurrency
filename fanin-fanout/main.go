package main

import (
	"fmt"
	"sync"
)

// getInputChan returns a channel for input numbers
func getInputChan() <-chan int {
	input := make(chan int, 100)
	nums := []int{1, 2, 3, 4, 5}

	go func() {
		defer close(input)
		for _, num := range nums {
			input <- num
		}
	}()

	return input
}

// getSquareChan returns a channel which returns square of input numbers
func getSquareChan(input <-chan int) <-chan int {
	output := make(chan int, 100)

	go func() {
		defer close(output)
		for num := range input {
			fmt.Println("getSquareChan: ", num)
			output <- num * num
		}
	}()

	return output
}

// merge returns a merged channel of outputsChan channels
func merge(outputsChan ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	merged := make(chan int, 100)
	wg.Add(len(outputsChan))

	// merged チャネルにまとめるための関数。引数のチャネルから数値を受け取って merged チャネルに送信する
	// ここでは getSquareChan によって二乗された数値を merged チャネルに送信するために使われている
	outputs := func(sc <-chan int) {
		defer wg.Done()
		for n := range sc {
			merged <- n
		}
	}

	// 上記 outpus 関数を受け取ったチャネルの数だけ起動する
	for _, optChan := range outputsChan {
		go outputs(optChan)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	inputChan := getInputChan()

	// fan-out: square operations to multiple goroutines
	// それぞれ inputChan から提供される数値の2乗を計算して送信するチャネルが返される
	// getSquareChan を2回呼び出しているので、同時実行数は2になる。順番は保証されない。
	squareChan1 := getSquareChan(inputChan)
	squareChan2 := getSquareChan(inputChan)

	// fan-in: squareChan1 と squareChan2 の結果を1つのチャネルにまとめる
	mergedChan := merge(squareChan1, squareChan2)

	for num := range mergedChan {
		fmt.Println(num)
	}
}
