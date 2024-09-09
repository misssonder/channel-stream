package channel_stream

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestForEachConcurrent(t *testing.T) {
	var (
		length = 100
		limit  = 10
		stream = make(chan int, length)
		result = make([]int, 0, length)
		lock   sync.Mutex
	)
	for i := 0; i < length; i++ {
		stream <- i
	}
	close(stream)
	ForEachConcurrent(stream, limit, func(i int) {
		time.Sleep(100 * time.Millisecond)
		lock.Lock()
		defer lock.Unlock()
		result = append(result, i)
	})
	sort.Ints(result)
	assert.Equal(t, length, len(result))
}

func TestBufferUnordered(t *testing.T) {
	var (
		length = 100
		limit  = 10
		stream = make(chan int64, length)
		result = make([]int, 0, length)
	)
	for i := 0; i < length; i++ {
		stream <- int64(i)
	}
	close(stream)
	outputStream := BufferUnordered(stream, limit, func(i int64) int {
		time.Sleep(100 * time.Millisecond)
		return int(i)
	})
	for i := range outputStream {
		result = append(result, i)
	}
	sort.Ints(result)
	assert.Equal(t, length, len(result))
	assert.Equal(t, func() []int {
		output := make([]int, 0, length)
		for i := 0; i < length; i++ {
			output = append(output, i)
		}
		return output
	}(), result)
}

func TestBuffered(t *testing.T) {
	var (
		length = 100
		limit  = 10
		stream = make(chan int64, length)
		result = make([]int, 0, length)
	)
	for i := 0; i < length; i++ {
		stream <- int64(i)
	}
	close(stream)
	outputStream := Buffered(stream, limit, func(i int64) int {
		time.Sleep(100 * time.Millisecond)
		return int(i)
	})
	for i := range outputStream {
		result = append(result, i)
	}
	// we don't sort the result, but it's still ordered
	//sort.Ints(result)
	assert.Equal(t, length, len(result))
	assert.Equal(t, func() []int {
		output := make([]int, 0, length)
		for i := 0; i < length; i++ {
			output = append(output, i)
		}
		return output
	}(), result)
}
