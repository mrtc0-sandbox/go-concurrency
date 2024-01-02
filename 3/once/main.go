package main

import (
	"fmt"
	"sync"
)

var i = 0

func warmup() {
	i++
}

func main() {
	var once sync.Once
	fmt.Println(i)
	once.Do(warmup)
	once.Do(warmup)
	fmt.Println(i)
}
