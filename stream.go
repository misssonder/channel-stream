package channel_stream

func Filter[T any](stream <-chan T, predicate func(val T) bool) <-chan T {
	result := make(chan T, cap(stream))

	go func() {
		defer close(result)
		for i := range stream {
			if predicate(i) {
				result <- i
			}
		}
	}()

	return result
}

func Map[T any, R any](stream <-chan T, iteratee func(T) R) <-chan R {
	result := make(chan R, cap(stream))

	go func() {
		defer close(result)
		for i := range stream {
			result <- iteratee(i)
		}
	}()

	return result
}

func FilterMap[T any, R any](stream <-chan T, callback func(item T) (R, bool)) <-chan R {
	result := make(chan R, cap(stream))

	go func() {
		defer close(result)
		for i := range stream {
			if r, ok := callback(i); ok {
				result <- r
			}
		}
	}()

	return result
}
