package main

import (
	"fmt"
)


type List struct {
	Value int
	Next *List
	Previous *List
}

func main() {
	var head *List
	var tail *List
	var n, val int
	fmt.Scanf("%d", &n)
	for i:=0; i<n; i++ {
		fmt.Scanf("%d", &val)
		list := List{}
		list.Value = val
		list.Next = nil
		list.Previous = nil
		//fmt.Println(list.Value)
		if i == 0 {
			head = &list
			tail = &list
		}

		if i > 0 {
			list.Previous = tail
			tail.Next = &list
			tail = &list
		}

	}

	for list := head; list !=nil; list = list.Next {
		fmt.Println(&list.Value,*list)
	}
	fmt.Println()
	fmt.Println()

	for list := tail; list != nil; list = list.Previous {
		fmt.Println(&list.Value,*list)
	}

}