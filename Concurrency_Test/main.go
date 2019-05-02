package main

import (
	"fmt"
	"sync"
	"time"
)
var wg sync.WaitGroup

func sqrt(jobs, results chan float64) {
	defer wg.Done()

	for val :=  range jobs {
		low := 0.0
		high := val
		mid := (low + high) / 2
		for {
			if high-low <= 0.000000001 {
				break
			}
			if mid*mid > val {
				high = mid
			} else {
				low = mid
			}
			mid = (low + high) / 2
		}
		time.Sleep(time.Nanosecond)
		results <- mid
	}
}

/*func sqrt(val float64) float64 {

		low := 0.0
		high := val
		mid := (low + high) / 2
		for {
			if high-low <= 0.000001 {
				break
			}
			if mid*mid > val {
				high = mid
			} else {
				low = mid
			}
			mid = (low + high) / 2
		}
		time.Sleep(time.Nanosecond)
		return mid

}*/

//666715.4592605451
//11.3265803

func main() {
	start := time.Now()
	sum := 0.0
	jobs := make(chan float64, 10000)
	results := make(chan float64, 10000)
	//var ret float64
	for w := 1; w<=1000; w++ {
		wg.Add(1)
		go sqrt(jobs, results)
	}

	for i:=2.0; i<=10000.0; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	close(results)

	for res := range results {
		sum += res
	}

	/*for i:=2.0; i<=10000.0; i++ {
		sum += sqrt(i)
	}*/

	fmt.Println(sum)
	fmt.Println(time.Since(start).Seconds())
}