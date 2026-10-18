package main

import (
	"runtime"
	"sync"
	"testing"
)

func FanOutMyWithChan[T any, K any](mas []T, f func(T) K, numberOfWork int) []<-chan K {

	channels := make([]<-chan K, numberOfWork)

	d := len(mas) / numberOfWork
	remainder := len(mas) % numberOfWork

	if remainder > 0 {
		numberOfWork--
	}

	work := func(ch chan K, start int, end int) {
		defer close(ch)
		for k := start; k < end; k++ {
			ch <- f(mas[k])
		}
	}

	for i := 0; i < numberOfWork; i++ {
		ch := make(chan K, 1000)
		channels[i] = ch
		go work(ch, i*d, (i+1)*d)
	}

	if remainder > 0 {
		ch := make(chan K, 1000)
		channels[numberOfWork] = ch
		go work(ch, numberOfWork*d, len(mas))
	}

	return channels
}

func FanInMyWithChan[T any](chs []<-chan T) []T {

	results := make([][]T, len(chs))

	var wg sync.WaitGroup

	for i, ch := range chs {
		iCopy := i
		chCopy := ch
		results[i] = make([]T, 0, 1000)

		wg.Go(func() {
			for value := range chCopy {
				results[iCopy] = append(results[iCopy], value)
			}
		})
	}

	wg.Wait()

	for i := 1; i < len(results); i++ {
		results[0] = append(results[0], results[i]...)
	}

	return results[0]
}

var testResultMyWithChan []int

func BenchmarkChannels(b *testing.B) {

	array := make([]int, 100_000_000)
	for i := range array {
		array[i] = i
	}

	b.ResetTimer()

	for b.Loop() {

		fanOut := FanOutMyWithChan(array, func(v int) int { return v * 2 }, runtime.GOMAXPROCS(0))
		testResultMyWithChan = FanInMyWithChan(fanOut)
	}
}
