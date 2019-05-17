package main

import (
	"errors"
	"fmt"
	"time"
)

type Stack struct {
	slice []int
	sleep time.Duration
}

func (s *Stack) Push(v int) {
	s.slice = append(s.slice, v)
	time.Sleep(time.Second * s.sleep)
}

func (s *Stack) Pop() (int, error) {
	if len(s.slice) == 0 {
		return 0, errors.New("Stack is empty")
	}
	ret := s.slice[len(s.slice)-1]
	s.slice = s.slice[0:len(s.slice)-1]
	time.Sleep(time.Second * s.sleep)
	return ret, nil
}

func (s *Stack) Peek() (int, error) {
	if len(s.slice) == 0 {
		return 0, errors.New("Stack is empty")
	}
	ret := s.slice[len(s.slice)-1]
	return ret, nil
}

func (s *Stack) String() string {
	return fmt.Sprint(s.slice)
}

func main() {
	start := time.Now()
	var s *Stack = new(Stack)
	s.sleep = 1
	fmt.Println(s.Peek())
	fmt.Println(s.Pop())
	s.Push(123)
	s.Push(234)
	s.Push(543)
	fmt.Println(s)
	fmt.Println(s.Peek())
	fmt.Println(s.Pop())
	fmt.Println(s)
	fmt.Println(s.Peek())
	//time.Sleep(time.Second*2)
	fmt.Println(s.Pop())
	fmt.Println(s)
	fmt.Println(s.Peek())
	//time.Sleep(time.Second*2)
	s.Push(600)
	fmt.Println(s)
	fmt.Println(s.Peek())
	fmt.Println(s.Pop())
	fmt.Println(s)
	fmt.Println(s.Peek())
	fmt.Println(s.Pop())
	//time.Sleep(time.Second*2)
	fmt.Println(s)
	fmt.Println(s.Pop())
	fmt.Println(time.Since(start).Seconds())
}