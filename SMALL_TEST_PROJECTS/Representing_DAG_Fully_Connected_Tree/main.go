package main

import "fmt"

type Node struct {
	Value int
	Left *Node
	Right *Node
}

func read() []Node {
	var n, val, indexLeft, indexRight int
	fmt.Scanf("%d", &n)
	var nodes = make([]Node, n)
	for i :=0; i<n; i++ {
		fmt.Scanf("%d %d %d", &val, &indexLeft, &indexRight)
		nodes[i].Value = val
		if indexLeft >= 0 {
			nodes[i].Left = &nodes[indexLeft]
		}
		if indexRight >= 0 {
			nodes[i].Right = &nodes[indexRight]
		}
	}
	return nodes
}

func main() {
	nodes := read()

	for i, _ := range nodes {
	 	fmt.Print(nodes[i].Value)
	 	if nodes[i].Left != nil {
			fmt.Print(" ", nodes[i].Left.Value)
		}else {
			fmt.Print("  ")
		}
	 	if nodes[i].Right != nil {
			fmt.Println(" ", nodes[i].Right.Value)
		}else {
			fmt.Println("  ")
		}
	}
}