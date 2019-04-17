package main

import (
	"fmt"
)

func main() {
	var arr = []string{"hello", "world", "what's", "up"}

	var parr []string
	parr = arr
	arr[0] = "Hello"
	fmt.Println(parr[0])
	fmt.Println(arr[0])
}
