package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func nilakanthaSequential(terms int) float64 {
	pi := 3.0
	for i := 1; i <= terms; i++ {
		pi += 4.0 * math.Pow(-1, float64(i+1)) / (float64(2*i) * float64(2*i+1) * float64(2*i+2))
	}
	return pi
}

func nilakanthaParallel(terms int, numWorkers int) float64 {
	var wg sync.WaitGroup
	var mu sync.Mutex
	pi := 3.0

	blockSize := terms / numWorkers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()

			start := w*blockSize + 1
			end := (w + 1) * blockSize
			if w == numWorkers-1 {
				end = terms
			}

			localPi := 0.0
			for i := start; i <= end; i++ {
				localPi += 4.0 * math.Pow(-1, float64(i+1)) / (float64(2*i) * float64(2*i+1) * float64(2*i+2))
			}

			mu.Lock()
			pi += localPi
			mu.Unlock()
		}(w)
	}

	wg.Wait()

	return pi
}

func measureTime(f func() float64) (time.Duration, float64) {
	var minTime time.Duration
	var result float64

	for i := 0; i < 10; i++ {
		start := time.Now()
		result = f()
		duration := time.Since(start)
		if minTime == 0 || duration < minTime {
			minTime = duration
		}
	}

	return minTime, result
}

func main() {
	terms := 10_000_000
	numWorkers := 8

	sequentialTime, sequentialPi := measureTime(func() float64 {
		return nilakanthaSequential(terms)
	})
	fmt.Printf("Sequential calculation took: %v, π ≈ %.15f\n", sequentialTime, sequentialPi)

	parallelTime, parallelPi := measureTime(func() float64 {
		return nilakanthaParallel(terms, numWorkers)
	})
	fmt.Printf("Parallel calculation took: %v, π ≈ %.15f\n", parallelTime, parallelPi)

	speedup := float64(sequentialTime) / float64(parallelTime)
	fmt.Printf("Speedup: %.2f\n", speedup)

	fmt.Printf("Math π for comparison: %.15f\n", math.Pi)
}
