# Golang local cache

Implementation of simple local cache.

[![Build Status](https://travis-ci.org/sima-land/go-local-cache.svg?branch=master)](https://travis-ci.org/sima-land/go-local-cache)
[![Go Report Card](https://goreportcard.com/badge/github.com/sima-land/go-local-cache)](https://goreportcard.com/report/github.com/sima-land/go-local-cache)

Features & todos:
- [x] prevents "dogpile effect"
- [x] check data expiration time
- [x] lru policy
- [ ] cache hit and enviction stats

Cache limitations:

- keys are always int
- get function returns multiple values 
- set function sets multiple keys at once
- to prevent dog pile effect queueKey should be manually set

## How to use

Create cache instance

```golang
import "github.com/sima-land/go-local-cache"
c := cache.New()
```

Create getter function. Its preferred way to fill the cache because in case of "dog pile" effect getter will be 
called once. 

```golang
getter := func(keys Keys) (Result, error) {
	r := Result{}
	for _, k := range keys {
		// haavy call to database
		r[k] = db.GetByID(k)		
	}
	return r, nil
}
```

Obtain data from cache. If there is no data in the cache getter will be called.

Note:
1. Next call is concurrently safe
2. If too many go-routines in our application are trying to call this getter with same queue key (user-1) 
   then only first call will actually be done

```golang
result, err := c.Get(int[]{1}, getter, "user-1")
```

## QueueKey param explain

You can simple use entity name and id as queueKey param of Get function. Eg:

```golang
result, err := c.Get(int[]{1}, userGetter, "user-1")
result, err := c.Get(int[]{1}, productGetter, "product-1")
```

Or if you need to get a page of product's block whole page. If keys 1,2 are in cache but other keys are not then:

1. Keys 1 and 2 will be get from cache
2. Getter will be called for keys 3,4,5 and 6
3. If another getter with same queue key "product-1" was called before and first one will be used instead

```golang
result, err := c.Get(int[]{1,2,3,4,5,6}, userGetter, "product-page-1")
```

## Options

Cache has a variety of configuration options:

``` golang
c := cache.New(cache.Options{
		MaxEntries: 1000000, // Specify max number items in cache
		TTL: time.Hour, // Time to live for cache entries
	})
```

## Test

MacBook 2,7 GHz Intel Core i5

```
BenchmarkCache_Get-4      500000              2274 ns/op
BenchmarkCache_Get-4      500000              2335 ns/op
```
