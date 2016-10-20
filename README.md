# Golang local cache

Implementation of simple local cache.

Features:
- prevents "dogpile effect" when too many go-routines in our application are trying to calculate new value to cache it.
- check data expiration time

