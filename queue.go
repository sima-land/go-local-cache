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

// Do place call with queue of calls with same key and return results of call.
// Actually only one getter will be called at the same time
func (q *Queue) Do(key interface{}, param Param, getter Getter) (Result, error) {
	q.mu.Lock()
	if q.calls == nil {
		q.calls = make(map[interface{}]*call)
	}
	if c, ok := q.calls[key]; ok {
		q.mu.Unlock()
		c.wg.Wait()
		return c.res, c.err
	}

	c := new(call)
	c.wg.Add(1)
	q.calls[key] = c
	q.mu.Unlock()
	c.res, c.err = getter(param)
	c.wg.Done()

	q.mu.Lock()
	delete(q.calls, key)
	q.mu.Unlock()

	return c.res, c.err
}