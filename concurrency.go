package channel_stream

import "sync"

type bufferedResult[T any] struct {
	result T
	index  int
}

func ForEachConcurrent[T any](stream <-chan T, limit int, callback func(T)) {
	var (
		concurrent = make(chan struct{}, limit)
		wg         = &sync.WaitGroup{}
	)

	for i := range stream {
		concurrent <- struct{}{}
		wg.Add(1)
		go func(val T) {
			defer func() {
				<-concurrent
				wg.Done()
			}()
			callback(val)
		}(i)
	}

	wg.Wait()
}

func BufferUnordered[T any, R any](stream <-chan T, bufferSize int, callback func(T) R) <-chan R {
	result := make(chan R, cap(stream))

	go func() {
		defer close(result)
		var (
			concurrent = make(chan struct{}, bufferSize)
			wg         = &sync.WaitGroup{}
		)

		for i := range stream {
			concurrent <- struct{}{}
			wg.Add(1)
			go func(val T) {
				defer func() {
					<-concurrent
					wg.Done()
				}()
				result <- callback(val)
			}(i)
		}
		wg.Wait()
	}()

	return result
}

func Buffered[T any, R any](stream <-chan T, bufferSize int, callback func(T) R) <-chan R {
	var (
		result          = make(chan R, cap(stream))
		wg              = &sync.WaitGroup{}
		unorderedResult = make(chan bufferedResult[R], bufferSize)
	)

	go func() {
		defer func() {
			wg.Wait()
			close(unorderedResult)
		}()
		var (
			index      = 0
			concurrent = make(chan struct{}, bufferSize)
		)
		// invoke callback with index
		for i := range stream {
			index += 1
			concurrent <- struct{}{}
			wg.Add(1)
			go func(index int, val T) {
				defer func() {
					<-concurrent
					wg.Done()
				}()
				unorderedResult <- bufferedResult[R]{
					result: callback(val),
					index:  index,
				}
			}(index, i)
		}
	}()

	// reorder results
	go func() {
		defer func() {
			close(result)
		}()
		var (
			reorder = make(map[int]R, bufferSize)
			expect  = 1
		)
		for i := range unorderedResult {
			reorder[i.index] = i.result
			for {
				val, ok := reorder[expect]
				if !ok {
					break
				}
				result <- val
				delete(reorder, expect)
				expect += 1
			}
		}
	}()
	return result
}
