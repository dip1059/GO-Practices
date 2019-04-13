package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var r io.Reader
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return
	}
	r = tty
	fmt.Println(r)
}
