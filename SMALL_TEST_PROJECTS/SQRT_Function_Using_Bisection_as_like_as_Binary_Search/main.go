package main

import (
	"fmt"
	"time"
)

func sqrt(x float64) float64 {
	low := 0.0
	high := x
	mid := (high + low) / 2
	for {
		if high - low <= 0.00000000001 {
			break
		}
		if mid * mid > x {
			high = mid
		} else {
			low = mid
		}
		mid = (high + low) / 2
	}
	return mid
}

func main() {
	start := time.Now()
	//fmt.Println(math.Sqrt(11000000.0))
	fmt.Println(sqrt(11000000.0))
	fmt.Println(time.Since(start).Nanoseconds()/10000)
}