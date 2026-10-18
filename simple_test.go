package main

import (
	"runtime"
	"sync"
	"testing"
)

func Array2Chan[T any](mas []T) <-chan T {
	result := make(chan T, len(mas))

	go func() {
		defer close(result)
		for i := range mas {
			result <- mas[i]
		}

	}()

	return result
}

func FanOutSimple[T any, K any](items <-chan T, f func(T) K, numberOfWork int) []<-chan K {

	result := make([]<-chan K, 0, numberOfWork)

	for range numberOfWork {
		res := make(chan K)
		result = append(result, res)
		go func() {
			defer close(res)
			for item := range items {
				res <- f(item)
			}
		}()
	}

	return result
}

func FanInSimple[T any](chs []<-chan T) <-chan T {

	result := make(chan T)

	var wg sync.WaitGroup

	for _, ch := range chs {
		chCopy := ch

		wg.Go(func() {
			for value := range chCopy {
				result <- value
			}
		})
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func Chan2Array[T any](ch <-chan T) []T {
	result := make([]T, 0, 16)
	for item := range ch {
		result = append(result, item)
	}

	return result
}

var testResultSimple []int

func BenchmarkChannel(b *testing.B) {

	array := make([]int, 100_000_000)
	for i := range array {
		array[i] = i
	}

	b.ResetTimer()

	for b.Loop() {

		chanItems := Array2Chan(array)
		fanOut := FanOutSimple(chanItems, func(v int) int { return v * 2 }, runtime.GOMAXPROCS(0))
		fanIn := FanInSimple(fanOut)
		testResultSimple = Chan2Array(fanIn)
	}
}

func BenchmarkSequential(b *testing.B) {

	array := make([]int, 100_000_000)
	for i := range array {
		array[i] = i
	}

	b.ResetTimer()

	for b.Loop() {

		result := make([]int, 100_000_000)
		for i := 0; i < 100_000_000; i++ {
			result[i] = array[i] * 2
		}

		testResultSimple = result
	}
}
