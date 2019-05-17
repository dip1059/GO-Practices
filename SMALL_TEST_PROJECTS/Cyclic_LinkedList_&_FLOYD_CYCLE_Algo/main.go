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
	var cycleStart *List
	var n, val int
	fmt.Scanf("%d", &n)
	for i:=0; i<n; i++ {
		fmt.Scanf("%d", &val)
		list := List{}
		list.Value = val
		list.Next = nil
		list.Previous = nil
		if i == 0 {
			head = &list
			tail = &list
		}

		//started the cycle from 2nd element
		if i == 1 {
			cycleStart = &list
		}
		//
		if i > 0 {
			list.Previous = tail
			tail.Next = &list
			tail = &list
		}

	}
	tail.Next = cycleStart

	hare := head
	tortoise := head
	tortoise2 := head

	//Floyd Cycle Detecting Algorithm
	for hare != nil {
		if hare.Next != nil && hare.Next.Next != nil {
			hare = hare.Next.Next
		} else {
			fmt.Println("Not Cyclic")
			break
		}

		tortoise = tortoise.Next

		if hare == tortoise {
			fmt.Println("Cyclic")
			break
		}
	}
	//

	//finding the Starting Node of the Cycle
	for {
		tortoise = tortoise.Next
		tortoise2 = tortoise2.Next
		if tortoise == tortoise2 {
			fmt.Println("Cycle starting Node", *tortoise)
			break
		}
	}
	//
}