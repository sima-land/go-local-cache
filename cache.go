package cache

import (
	"sync"
)

// Getter func must be implemented by client and returns data by it's keys
type Getter func(keys Keys) (Result, error)

// Result represents result returned by cache function
type Result map[int]interface{}

// Keys represents set of cache keys
// It always slice of int for two reasons:
// - we do not have separate Get and GetMulti functions
// - integers are most common used primary keys in db
type Keys []int

// Cache represents cache
type Cache struct {
	getter Getter
	data   Result
	mu     sync.RWMutex
	queue  *Queue
}

// New initializes new instance of cache
func New() *Cache {
	c := Cache{}
	c.data = make(Result)
	c.queue = NewQueue()
	return &c;
}


// Get returns data from cache by keys
//
// If data not exist Getter function will be called to get data
// Data returned by Getter function will be placed into cache
//
func (c *Cache) Get(keys Keys, getter Getter, queueKey interface{}) (Result, error) {
	result := make(Result)
	missedKeys := make(Keys, 0)
	c.mu.RLock()
	for _, key := range keys {
		val, ok := c.data[key]
		if ok {
			result[key] = val
		} else {
			missedKeys = append(missedKeys, key)
		}
	}
	c.mu.RUnlock()
	if len(missedKeys) == 0 {
		return result, nil
	}
	missedResult, err := c.queue.Do(missedKeys, getter, queueKey)
	if err != nil {
		return result, err
	}
	c.Set(missedResult)
	for k, d := range missedResult {
		result[k] = d
	}
	return result, nil
}

// Set saves result into cache
func (c *Cache) Set(data Result) {
	c.mu.Lock()
	for k, d := range data {
		c.data[k] = d
	}
	c.mu.Unlock()
}

// Len returns the number of items in the cache
func (c *Cache) Len() int {
	if c.data == nil {
		return 0
	}
	return len(c.data)
}
