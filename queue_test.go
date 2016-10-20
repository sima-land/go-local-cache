package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueue_Do(t *testing.T) {
	q := NewQueue()
	var actual1, actual2 Result
	sample1 := Result{1 : true}
	sample2 := Result{1 : false}

	results := make(chan Result, 2)
	finish := make(chan bool)
	start := make(chan bool)

	getter := func(param Param) (Result, error) {
		result := <-results
		return result, nil
	}

	go func() {
		start <- true
		actual1, _ = q.Do(1, nil, getter)
		finish <- true
	}()
	<-start
	go func() {
		start <- true
		actual2, _ = q.Do(1, nil, getter)
		finish <- true
	}()
	<-start

	results <- sample1
	results <- sample2

	<-finish
	<-finish

	//require.Equal(t, 1, int(counter))

	// Make sure getter was called once
	// And second getter call sample of first one
	require.Equal(t, sample1, actual1)
	require.Equal(t, sample1, actual2)

	// Make sure getter called one
	require.Equal(t, 1, len(results))
}