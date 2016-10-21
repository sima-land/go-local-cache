package cache

import (
	"sync"
)

// NewQueue initialize new Queue
func NewQueue() *Queue {
	q := Queue{}
	q.calls = make(map[interface{}]*call)
	return &q
}

// Queue of calls. Allows to prevent go pile effect
type Queue struct {
	mu    sync.Mutex
	calls map[interface{}]*call
}

// call information about single call
type call struct {
	wg  sync.WaitGroup
	res Result
	err error
}

// Do place call in queue of calls with same queueKey
// and returns results of first call in queue
func (q *Queue) Do(keys Keys, getter Getter, queueKey interface{}) (Result, error) {
	q.mu.Lock()
	if q.calls == nil {
		q.calls = make(map[interface{}]*call)
	}
	if c, ok := q.calls[queueKey]; ok {
		q.mu.Unlock()
		c.wg.Wait()
		return c.res, c.err
	}

	c := new(call)
	c.wg.Add(1)
	q.calls[queueKey] = c
	q.mu.Unlock()
	c.res, c.err = getter(keys)
	c.wg.Done()

	q.mu.Lock()
	delete(q.calls, queueKey)
	q.mu.Unlock()

	return c.res, c.err
}