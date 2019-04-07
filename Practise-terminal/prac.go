package main

import (
	"fmt"
)

func main() {
	//var wg sync.WaitGroup
	c1 := make(chan string)
	c2 := make(chan string)
	//wg.Add(2)
		go PrintWord(c1)
		msg := <- c1
		fmt.Println(msg)
		//wg.Done()
	go	PrintNumber(c2)
		msg = <- c2
		fmt.Println(msg)
		//wg.Done()
	//wg.Wait()
}

func PrintWord(c1 chan string) {
	c1 <- "This is PrintWord"
	//close(c1)
}

func PrintNumber(c2 chan string) {
	c2 <- "This is PrintNumber"
	close(c2)
}
