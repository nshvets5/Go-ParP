package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func generateMatrix(n int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = rand.Intn(100)
		}
	}
	return matrix
}

func matrixMultiplySequential(A, B [][]int, n int) [][]int {
	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, n)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return C
}

func matrixMultiplyParallel(A, B [][]int, n int) [][]int {
	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, n)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				for k := 0; k < n; k++ {
					C[i][j] += A[i][k] * B[k][j]
				}
			}(i, j)
		}
	}

	wg.Wait()
	return C
}

func measureTime(f func() [][]int) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func main() {
	n := 500

	A := generateMatrix(n)
	B := generateMatrix(n)

	sequentialTime := measureTime(func() [][]int {
		return matrixMultiplySequential(A, B, n)
	})
	fmt.Printf("Sequential matrix multiplication took: %v\n", sequentialTime)

	parallelTime := measureTime(func() [][]int {
		return matrixMultiplyParallel(A, B, n)
	})
	fmt.Printf("Parallel matrix multiplication took: %v\n", parallelTime)

	speedup := float64(sequentialTime) / float64(parallelTime)
	fmt.Printf("Speedup: %.2f\n", speedup)
}
