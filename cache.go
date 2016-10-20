package cache

// Getter func must be implemented by client
type Getter func(param Param) (Result, error)

// Result represents result returned by cache function
type Result map[int]interface{}

// Param represents param to get cache data
// It always slice of int because we do not have separate Get and GetMulti functions.
// Get function always returns many cached objects
type Param []int