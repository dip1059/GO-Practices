package main

import (
	"fmt"
	"time"
)

func BinarySearch(arr []int, val int) func() (int, bool, bool) {
	low := 0
	high := len(arr) - 1
	mid := (low + high) / 2
	res := false
	finish := false

	return func() (int, bool, bool){
		if low > high {
			finish = true
			return mid, res, finish
		} else {
			if val > arr[mid] {
				low = mid + 1
			} else if val < arr[mid] {
				high = mid - 1
			} else {
				res = true
				finish = true
				return mid, res, finish
			}
			mid = (low + high) / 2
		}
		return mid, res, finish
	}
}

func main() {
	start := time.Now()
	arr := []int{}
	for i := 12; true ; i += time.Now().Second() {
		arr = append(arr, i)
		if len(arr) == 1000 {
			break
		}
	}

	for _, v := range arr{
		fmt.Print(v," ")
	}
	fmt.Println()
	val := 876
	/*low := 0
	high := len(arr) - 1
	mid := (low + high) / 2
	res := false
	for {
		if low > high {
			break
		} else {
			if val > arr[mid] {
				low = mid + 1
			} else if val < arr[mid] {
				high = mid - 1
			} else if val == arr[mid] {
				res = true
				break
			}
			mid = (low + high) / 2
		}
	}*/
	BS := BinarySearch(arr, val)
	for {
		mid, res, finish := BS()
		if finish {
			if res {
				fmt.Println(val, "Found at index", mid)
				fmt.Println(time.Since(start).Nanoseconds()/1000, "microsecond")
			} else {
				fmt.Println(val, "Not found.")
				fmt.Println(time.Since(start).Nanoseconds()/1000, "microsecond")
			}
			break
		}
	}

}
