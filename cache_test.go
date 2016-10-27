package cache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func getter(keys Keys) (Result, error) {
	r := Result{}
	for _, k := range keys {
		r[k] = k
	}
	return r, nil
}

func TestCache_Get(t *testing.T) {

	c := New(Options{MaxEntries: 4})
	c.Set(Result{1: 1, 2: 2})

	r, _ := c.Get([]int{1, 2}, getter, 1)
	require.Equal(t, 1, r[1])
	require.Equal(t, 2, r[2])
	require.Equal(t, 2, c.Len())

	r, _ = c.Get([]int{3, 4}, getter, 1)
	require.Equal(t, 3, r[3])
	require.Equal(t, 4, r[4])
	require.Equal(t, 4, c.Len())

	r, _ = c.Get([]int{5, 6}, getter, 1)
	require.Equal(t, 5, r[5])
	require.Equal(t, 6, r[6])

	require.Equal(t, 4, c.Len())
}

func TestCache_TTL(t *testing.T) {
	ttl := time.Millisecond * 10

	c := New(Options{TTL:ttl})
	c.Set(Result{1: 10, 2: 20})
	r, _ := c.Get([]int{1, 2}, getter, 1)
	require.Equal(t, 10, r[1])
	require.Equal(t, 20, r[2])

	time.Sleep(ttl)
	r, _ = c.Get([]int{1, 2}, getter, 1)
	require.Equal(t, 1, r[1])
	require.Equal(t, 2, r[2])
}

func BenchmarkCache_Get(b *testing.B) {

	const maxPage = 10
	const PageSize = 10

	getter := func(keys Keys) (Result, error) {
		r := Result{}
		for _, k := range keys {
			r[k] = k
		}
		time.Sleep(time.Millisecond * 10)
		return r, nil
	}

	c := New(Options{MaxEntries: maxPage*PageSize - PageSize})

	// test func request data from cache
	test := func() {
		start := int(rand.Int31n(maxPage-1) + 1) // 1...MaxPage
		for i := start; i <= start+10; i++ {
			page := i
			if i > maxPage {
				page = i - maxPage
			}
			keys := make([]int, PageSize, PageSize)
			for j := 0; j < PageSize; j++ {
				keys[j] = page*PageSize + j
			}
			c.Get(Keys(keys), getter, page)
		}
	}

	for n := 0; n < b.N; n++ {
		go test()
	}
}
