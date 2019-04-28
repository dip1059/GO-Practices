package main

import (
	//custom package
	"Practice-Goland/algo"
	//
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	arr := []int{}
	for i := 12; true ; i += 2 {
		arr = append(arr, i)
		lenArr := len(arr)
		fmt.Print(arr[lenArr-1]," ")
		if lenArr == 1000 {
			break
		}
	}
	fmt.Println()

	val := 2010
	low := 0
	high := len(arr) - 1
	mid := (low + high) / 2

	//custom function BinarySearch()
	index, res := algo.BinarySearch(arr, val, low, high, mid)
	//
	//index, res := algo.LinearSearch(arr, val)

	if res {
		fmt.Println(val, "Found at index", index)
		fmt.Println(time.Since(start).Nanoseconds()/1000000, "millisecond")
	} else {
		fmt.Println(val, "Not found.")
		fmt.Println(time.Since(start).Nanoseconds()/1000000, "millisecond")
	}
}
