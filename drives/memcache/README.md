# MemCache


## Installation

```shell
go get -u github.com/gojuukaze/YTask/drives/memcache/v3
```

## Backend

```go
package main

import (
    "github.com/gojuukaze/YTask/drives/memcache/v3"
)

func main() {
	backend := memcache.NewMemCacheBackend([]string{"127.0.0.1:11211"}, 10)
	// ...
}
```
