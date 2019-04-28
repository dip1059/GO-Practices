package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	a := "9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"
	b := "1"
	var sl []byte
	var sl2 []byte
	r := 0
	d := 0
	lenA := len(a)
	lenB := len(b)
	lenS := 0
	lenL := 0
	sum := 0
	c := 0
	if lenA > lenB {
		lenL = lenA
		lenS = lenB
	}else {
		lenL = lenB
		lenS = lenA
	}
	j := 0
	for i :=lenL-1; i>=0; i-- {
		j++
		if i>=(lenL-lenS) {
			sum = c + int((a[lenA-j] - 48)+(b[lenB-j]-48))
		}else {
			if lenL == lenA {
				sum = c + int(a[lenA-j] - 48)
			}else {
				sum = c + int(b[lenB-j] - 48)
			}
		}
		d = sum
		for {
			r = d % 10
			sl = append(sl, byte(r))
			d = d / 10
			c = d
			if i>0 || d == 0 {
				break
			}
		}
	}

	for i:=len(sl)-1; i>=0; i-- {
		sl2 = append(sl2, sl[i])
	}

	for _,v := range sl2 {
		fmt.Print(v)
	}
	fmt.Println("\n",time.Since(start).Nanoseconds()/1000,"microsecond")

}