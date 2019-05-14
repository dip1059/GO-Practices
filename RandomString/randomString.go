package main

import (
	"math/rand"
	"time"
)

var Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
var lenLetters = len(Letters)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	//sec := time.Now().Second()
	b := make([]byte, n)
	for i := range b {
		//b[i] = Letters[rand.Intn(sec)]
		b[i] = Letters[rand.Intn(lenLetters)]
	}
	return string(b)
}