package cache

import (
	"container/list"
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

type entry struct {
	key   int
	value interface{}
}

// Options keeps cache options
type Options struct {
	MaxEntries int
}

// Cache represents cache
type Cache struct {
	getter Getter
	data   map[int]*list.Element
	// opt keeps cache options
	opt   Options
	ll    *list.List
	mu    sync.RWMutex
	queue *Queue
}

// New initializes new instance of cache
func New(options ...Options) *Cache {
	var o Options
	if len(options) == 0 {
		o = Options{}
	} else {
		o = options[0]
	}
	c := Cache{opt: o}
	c.data = make(map[int]*list.Element)
	c.ll = list.New()
	c.queue = NewQueue()
	return &c
}

// Get returns data from cache by keys
//
// If data not exist Getter function will be called to get data
// Data returned by Getter function will be placed into cache
func (c *Cache) Get(keys Keys, getter Getter, queueKey interface{}) (Result, error) {
	result := make(Result)
	missedKeys := make(Keys, 0)
	c.mu.Lock()
	for _, key := range keys {
		if e, ok := c.data[key]; ok {
			c.ll.MoveToFront(e)
			result[key] = e.Value.(*entry).value
		} else {
			missedKeys = append(missedKeys, key)
		}
	}
	c.mu.Unlock()
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
	for k, v := range data {
		if e, ok := c.data[k]; ok {
			c.ll.MoveToFront(e)
			e.Value.(*entry).value = v
		} else {
			c.data[k] = c.ll.PushFront(&entry{key: k, value: v})
		}
	}
	c.mu.Unlock()
	if c.opt.MaxEntries != 0 && c.ll.Len() > c.opt.MaxEntries {
		c.Flush(c.ll.Len() - c.opt.MaxEntries)
	}
}

// Len returns the number of items in the cache
func (c *Cache) Len() int {
	return c.ll.Len()
}

// Remove removes the provided keys from the cache.
func (c *Cache) Remove(keys Keys) {
	c.mu.Lock()
	for _, k := range keys {
		if e, ok := c.data[k]; ok {
			c.ll.Remove(e)
			delete(c.data, k)
		}
	}
	c.mu.Unlock()
}

// Flush removes last (rare used) n elements from cache
func (c *Cache) Flush(n int) {
	c.mu.Lock()
	for i := 1; i <= n; i++ {
		if e := c.ll.Back(); e != nil {
			c.ll.Remove(e)
			delete(c.data, e.Value.(*entry).key)
		}
	}
	c.mu.Unlock()
}
