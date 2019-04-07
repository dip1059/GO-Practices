package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	c1 := make(chan string)
	c2 := make(chan string)
	wg.Add(2)
	go PrintWord(c1)
	msg1 := <-c1
	fmt.Println(msg1)
	wg.Done()
	go PrintNumber(c2)
	msg := <-c2
	fmt.Println(msg)
	wg.Done()
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

func PrintWord(c1 chan string) {
	c1 <- "This is PrintWord"
	close(c1)
}

func PrintNumber(c2 chan string) {
	c2 <- "This is PrintNumber"
	close(c2)
}
