package main

import (
	"fmt"
)

func main() {
	a := []string{"A", "B", "C", "D", "E"}
i := 2

// Remove the element at index i from a.
copy(a[i:], a[i+1:])
fmt.Println(a) // Shift a[i+1:] left one index.
a[len(a)-1] = ""
fmt.Println(a)
     // Erase last element (write zero value).
a = a[:len(a)-1]     // Truncate slice.

fmt.Println(a)
}
