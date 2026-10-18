package main

import (
	"fan_in__fan_out/unbounded"
	"runtime"
	"sync"
	"testing"
)

func FanOutUnbounded[T any, K any](mas []T, f func(T) K, numberOfWork int) []*unbounded.UnboundedChannel[K] {

	channels := make([]*unbounded.UnboundedChannel[K], numberOfWork)

	d := len(mas) / numberOfWork
	remainder := len(mas) % numberOfWork

	if remainder > 0 {
		numberOfWork--
	}

	work := func(ch *unbounded.UnboundedChannel[K], start int, end int) {
		defer ch.Close()
		for k := start; k < end; k++ {
			ch.Send(f(mas[k]))
		}
	}

	for i := 0; i < numberOfWork; i++ {
		ch := unbounded.New[K]()
		channels[i] = ch
		go work(ch, i*d, (i+1)*d)
	}

	if remainder > 0 {
		ch := unbounded.New[K]()
		channels[numberOfWork] = ch
		go work(ch, numberOfWork*d, len(mas))
	}

	return channels
}

func FanInUnbounded[T any](chs []*unbounded.UnboundedChannel[T]) []T {

	results := make([][]T, len(chs))

	var wg sync.WaitGroup

	for i, ch := range chs {
		iCopy := i
		chCopy := ch
		results[i] = make([]T, 0, 1000)

		wg.Go(func() {

			for {
				value, ok := chCopy.Receive()
				if !ok {
					break
				}

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

var testResultUnbounded []int

func BenchmarkUnbounded(b *testing.B) {

	array := make([]int, 100_000_000)
	for i := range array {
		array[i] = i
	}

	b.ResetTimer()

	for b.Loop() {

		fanOut := FanOutUnbounded(array, func(v int) int { return v * 2 }, runtime.GOMAXPROCS(0))
		testResultUnbounded = FanInUnbounded(fanOut)
	}
}
