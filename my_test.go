package main

import (
	"runtime"
	"sync"
	"testing"
)

func FanMy[T any, K any](mas []T, f func(T) K, numberOfWork int) []K {

	results := make([]*[]K, numberOfWork)

	d := len(mas) / numberOfWork
	remainder := len(mas) % numberOfWork

	if remainder > 0 {
		numberOfWork--
	}

	var wg sync.WaitGroup

	work := func(res *[]K, start int, count int) {
		defer wg.Done()

		for k := range count {
			*res = append(*res, f(mas[start+k]))
		}
	}

	for i := range numberOfWork {
		wg.Add(1)
		start := i * d
		end := (i + 1) * d
		count := end - start
		res := make([]K, 0, count)
		results[i] = &res

		go work(results[i], start, count)
	}

	if remainder > 0 {
		wg.Add(1)
		start := numberOfWork * d
		end := len(mas)
		count := end - start
		res := make([]K, 0, count)
		results[numberOfWork] = &res

		go work(results[numberOfWork], start, count)
	}

	wg.Wait()

	for i := 1; i < len(results); i++ {
		*results[0] = append(*results[0], *results[i]...)
	}

	return *results[0]
}

var testResultMy []int

func BenchmarkArray(b *testing.B) {

	array := make([]int, 100_000_000)
	for i := range array {
		array[i] = i
	}

	b.ResetTimer()

	for b.Loop() {
		testResultMy = FanMy(array, func(v int) int { return v * 2 }, runtime.GOMAXPROCS(0))
	}
}
