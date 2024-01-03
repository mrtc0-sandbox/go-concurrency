package main

import (
	"fmt"
	"time"
)

// 複数のチャネルのうち、どれかが閉じられたら、まとめて閉じるサンプル
func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	// 複数のチャネルを受け取って、1つのチャネルを返す再帰関数
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		// 再帰関数として呼ばれるので、その終了条件。
		// チャネルが1つしかない場合は、そのチャネルを返す。そもそもチャネルがない場合は nil を返す
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				// チャネルが2つしかない場合の特別な処理
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
