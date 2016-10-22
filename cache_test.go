package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCache_Get(t *testing.T) {
	getter := func(keys Keys) (Result, error) {
		r := Result{}
		for _, k := range keys {
			r[k] = k
		}
		return r, nil
	}

	c := New()
	c.Set(Result{1:1, 2:2})

	r, _ := c.Get([]int{1,2}, getter, 1)
	require.Equal(t, 1, r[1])
	require.Equal(t, 2, r[2])
	require.Equal(t, 2, c.Len())

	r, _ = c.Get([]int{3,4}, getter, 1)
	require.Equal(t, 3, r[3])
	require.Equal(t, 4, r[4])
	require.Equal(t, 4, c.Len())
}