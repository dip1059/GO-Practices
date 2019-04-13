package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	r := rand.Intn(100)
	for {

		fmt.Printf("%d %.2fs\n", r, time.Since(start).Seconds())
	}
	//var st2 string = "1000"
	fmt.Printf("%d %.2fs\n", r, time.Since(start).Seconds())

}
