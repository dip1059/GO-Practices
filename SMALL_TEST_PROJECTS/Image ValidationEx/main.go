package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {

	// maximize CPU usage for maximum performance
	log.Println("cpu:",runtime.GOMAXPROCS(runtime.NumCPU()))

	// open the uploaded file
	file, err := os.Open("./s.jpg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buff := make([]byte, 512)
	log.Println("buff:", string(buff))// why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	i, err := file.Read(buff)
	log.Println("int", i)
	log.Println("buff:", string(buff))

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	fmt.Println(filetype)

	switch filetype {
	case "image/jpeg", "image/jpg":
		fmt.Println(filetype)

	case "image/gif":
		fmt.Println(filetype)

	case "image/png":
		fmt.Println(filetype)

	case "application/pdf": // not image, but application !
		fmt.Println(filetype)
	default:
		fmt.Println("unknown file type uploaded")
	}
}