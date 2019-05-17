package main

import "fmt"

type rect struct {
	h int
	w int
}

func (r *rect) area() (string, int) {
	return "Area is:", r.h * r.w
}

func main() {
	x := rect{h: 100, w: 5}
	var i int
	i = 12
	p := *rect{}
	/*x.h = 10
	x.w = 5*/
	fmt.Println(x.area())
	fmt.Println(i)
	fmt.Println(p)
}
