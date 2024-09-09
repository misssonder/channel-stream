package channel_stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var (
		length = 100
		stream = make(chan int, length)
	)
	for i := 0; i < length; i++ {
		stream <- i
	}
	close(stream)
	outputStream := Filter(stream, func(i int) bool {
		return i%2 == 0
	})
	for i := range outputStream {
		assert.True(t, i%2 == 0)
	}
}

func TestMap(t *testing.T) {
	var (
		length = 100
		stream = make(chan int, length)
	)
	for i := 0; i < length; i++ {
		stream <- i
	}
	close(stream)
	outputStream := Map(stream, func(i int) int64 {
		return int64(i)
	})
	expect := int64(0)
	for i := range outputStream {
		assert.Equal(t, expect, i)
		expect += 1
	}
}

func TestFilterMap(t *testing.T) {
	var (
		length = 100
		stream = make(chan int, length)
	)
	for i := 0; i < length; i++ {
		stream <- i
	}
	close(stream)
	outputStream := FilterMap(stream, func(i int) (int64, bool) {
		return int64(i), i%2 == 0
	})
	expect := int64(0)
	for i := range outputStream {
		assert.True(t, i%2 == 0)
		assert.Equal(t, expect, i)
		expect += 2
	}
}
