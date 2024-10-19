![ci](https://github.com/misssonder/channel-stream/actions/workflows/go.yml/badge.svg)
# Channel Stream
`channel-stream` is a Golang library provides functions to handle `chan` gracefully.

# Install
```shell
go get github.com/misssonder/channel-stream
```

# Example
## Filter
```go
var (		
	length = 100
	stream = make(chan int, length)
)
for i := 0; i < length; i++ {		
	stream <- i
}	
close(stream)	
outputStream := channle_stream.Filter(stream, func(i int) bool {		
	return i%2 == 0
})	
for i := range outputStream {		
	assert.True(t, i%2 == 0)
}
```

## Map
```go
var (		
	length = 100
	stream = make(chan int, length)
)	
for i := 0; i < length; i++ {		
	stream <- i
}	
close(stream)	
outputStream := channle_stream.Map(stream, func(i int) int64 {		
	return int64(i)	
})	
expect := int64(0)	
for i := range outputStream {		
	assert.Equal(t, expect, i)
	expect += 1
}
```

## FilterMap
```go
var (
	length = 100
	stream = make(chan int, length)
)
for i := 0; i < length; i++ {
	stream <- i
}
close(stream)
outputStream := channle_stream.FilterMap(stream, func(i int) (int64, bool) {
	return int64(i), i%2 == 0
})
expect := int64(0)
for i := range outputStream {
	assert.True(t, i%2 == 0)
	assert.Equal(t, expect, i)
	expect += 2
}
```
## ForEachConcurrent
like [for_each_concurrent](https://docs.rs/futures-util/latest/futures_util/stream/trait.StreamExt.html#method.for_each_concurrent) in rust library `futures_util`.
```go
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
channel_stream.ForEachConcurrent(stream, limit, func(i int) {
	time.Sleep(100 * time.Millisecond)
	lock.Lock()
	defer lock.Unlock()
	result = append(result, i)
})
sort.Ints(result)
assert.Equal(t, length, len(result))
```
## BufferUnordered
like [buffer_unordered](https://docs.rs/futures-util/latest/futures_util/stream/trait.StreamExt.html#method.buffer_unordered) in rust library `futures_util`. It will call the callback function asynchronously and it doesn't guarantee the order of output.
```go
var (
    length = 100
    limit = 10
    stream = make(chan int64, length)
    result = make([]int, 0, length)
)
for i := 0; i < length; i++ {
    stream <- int64(i)
}
close(stream)
outputStream := channel_stream.BufferUnordered(stream, limit, func (i int64) int {
time.Sleep(100 * time.Millisecond)
    return int(i)
})
for i := range outputStream {
    result = append(result, i)
}
sort.Ints(result)
assert.Equal(t, length, len(result))
assert.Equal(t, func () []int {
output := make([]int, 0, length)
for i := 0; i < length; i++ {
    output = append(output, i)
}
return output
}(), result)
```
## Buffered
like [buffered](https://docs.rs/futures-util/latest/futures_util/stream/trait.StreamExt.html#method.buffered) in rust library `futures_util`. It will call the callback function asynchronously and guarantee the order of result.
```go
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
outputStream := channel_stream.Buffered(stream, limit, func(i int64) int {
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
```