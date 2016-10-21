package cache

// Getter func must be implemented by client and returns data by it's keys
type Getter func(param Keys) (Result, error)

// Result represents result returned by cache function
type Result map[int]interface{}

// Keys represents set of cache keys
// It always slice of int for two reasons:
// - we do not have separate Get and GetMulti functions
// - integers are most common used primary keys in db
type Keys []int
