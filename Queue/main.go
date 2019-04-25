package main

import (
	"errors"
	"fmt"
	"time"
)

type Queue struct {
	slice []int
	sleep time.Duration
}

func (q *Queue) Enqueue(v int) {
	q.slice = append(q.slice, v)
	time.Sleep(time.Second * q.sleep)
}

func (q *Queue) Dequeue() (int, error) {
	if len(q.slice) == 0 {
		return 0, errors.New("Queue is empty.")
	}
	ret := q.slice[0]
	q.slice = q.slice[1:len(q.slice)]
	time.Sleep(time.Second * q.sleep)
	return ret, nil
}

func (q *Queue) String() string {
	return fmt.Sprint(q.slice)
}

func main() {
	start := time.Now()
	var q *Queue = new(Queue)
	q.sleep = 1
	fmt.Println(q.Dequeue())
	q.Enqueue(123)
	q.Enqueue(234)
	q.Enqueue(543)
	fmt.Println(q)
	fmt.Println(q.Dequeue())
	//time.Sleep(time.Second*2)
	fmt.Println(q.Dequeue())
	//time.Sleep(time.Second*2)
	q.Enqueue(600)
	fmt.Println(q)
	fmt.Println(q.Dequeue())
	//time.Sleep(time.Second*2)
	fmt.Println(q)
	fmt.Println(q.Dequeue())
	fmt.Println(q)
	fmt.Println(q.Dequeue())
	fmt.Println(time.Since(start).Seconds())
}