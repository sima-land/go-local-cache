# Golang local cache

Implementation of simple local cache.

[![Build Status](https://travis-ci.org/sima-land/go-local-cache.svg?branch=master)](https://travis-ci.org/sima-land/go-local-cache)
[![Go Report Card](https://goreportcard.com/badge/github.com/sima-land/go-local-cache)](https://goreportcard.com/report/github.com/sima-land/go-local-cache)

Features & todos:
- [x] prevents "dogpile effect"
- [ ] check data expiration time
- [ ] lru policy

## How to use

```golang

import "github.com/sima-land/go-local-cache"

getter := func(keys Keys) (Result, error) {
	r := Result{}
	for _, k := range keys {
		// obtain data for database
	}
	return r, nil
}

c := cache.New()

// 1. Next call is concurrently safe
// 2. If too many go-routines in our application are trying to call
// this getter with queue key equals 2 then only first call will actually be done
result, err := c.Get(int[]{1,2,3,4}, getter, 2)

```


