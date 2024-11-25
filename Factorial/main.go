package main

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"
)

func factorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

func measureTime(f func()) time.Duration {
	var minTime time.Duration
	for i := 0; i < 10; i++ {
		start := time.Now()
		f()
		duration := time.Since(start)
		if minTime == 0 || duration < minTime {
			minTime = duration
		}
	}
	return minTime
}

func sequentialFactorials(nums []int) {
	for _, n := range nums {
		_ = factorial(n)
	}
}

func parallelFactorials(nums []int) {
	var wg sync.WaitGroup
	ch := make(chan *big.Int, len(nums))

	for _, n := range nums {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ch <- factorial(n)
		}(n)
	}

	wg.Wait()
	close(ch)

	for result := range ch {
		_ = result
	}
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	nums := []int{5000, 10000, 15000, 20000, 25000}

	sequentialTime := measureTime(func() {
		sequentialFactorials(nums)
	})
	fmt.Printf("Sequential execution took: %v\n", sequentialTime)

	parallelTime := measureTime(func() {
		parallelFactorials(nums)
	})
	fmt.Printf("Parallel execution took: %v\n", parallelTime)

	speedup := float64(sequentialTime) / float64(parallelTime)
	fmt.Printf("Speedup: %.2f\n", speedup)
}
