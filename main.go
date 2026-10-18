package main

// func FanOut[T any, K any](mas []T, f func(T) K, numberOfWork int) []<-chan K {

// 	channels := make([]<-chan K, numberOfWork)

// 	d := len(mas) / numberOfWork
// 	remainder := len(mas) % numberOfWork

// 	if remainder > 0 {
// 		numberOfWork--
// 	}

// 	work := func(ch chan K, start int, end int) {
// 		defer close(ch)
// 		for k := start; k < end; k++ {
// 			ch <- f(mas[k])
// 		}
// 	}

// 	for i := 0; i < numberOfWork; i++ {
// 		ch := make(chan K)
// 		channels[i] = ch
// 		go work(ch, i*d, (i+1)*d)
// 	}

// 	if remainder > 0 {
// 		ch := make(chan K)
// 		channels[numberOfWork] = ch
// 		go work(ch, numberOfWork*d, len(mas))
// 	}

// 	return channels

// 	// result := make([]<-chan K, 0, numberOfWork)
// 	// items := make(chan T, len(mas))
// 	// defer close(items)

// 	// for range numberOfWork {
// 	// 	res := make(chan K, 10000)
// 	// 	result = append(result, res)
// 	// 	go func() {
// 	// 		defer close(res)
// 		for item := range items {
// 			res <- f(item)
// 		}
// 	}()
// }

// for i := range mas {
// 	items <- mas[i]
// }

// return result
// }

// func FanIn[T any](chs []<-chan T) <-chan T {

// 	result := make(chan T, 10000)

// 	var wg sync.WaitGroup

// 	for _, ch := range chs {
// 		chCopy := ch

// 		wg.Go(func() {
// 			for value := range chCopy {
// 				result <- value
// 			}
// 		})
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(result)
// 	}()

// 	return result
// }

// func ToArrayFromChanel[T any](ch <-chan T) []T {
// 	result := make([]T, 0, 16)
// 	for item := range ch {
// 		result = append(result, item)
// 	}

// 	return result
// }

// func ToArrayFromChannels[T any](chs []<-chan T) []T {
// 	results := make([][]T, len(chs))

// 	var wg sync.WaitGroup

// 	for i, ch := range chs {
// 		iCopy := i
// 		results[i] = make([]T, 0, 1000)

// 		wg.Go(func() {

// 			for value := range ch {
// 				results[iCopy] = append(results[iCopy], value)
// 			}
// 		})
// 	}

// 	wg.Wait()

// 	for i := 1; i < len(results); i++ {
// 		results[0] = append(results[0], results[i]...)
// 	}

// 	return results[0]
// }

func main() {

	// array := make([]int, 100_000)
	// for i := range array {
	// 	array[i] = i
	// }

	// res := ToArrayFromChannels(FanOut(array, func(v int) int { return v * 2 }, runtime.GOMAXPROCS(0)))
	// fmt.Println(len(res))

}
