# Redis 

## Installation

```shell
go get -u github.com/gojuukaze/YTask/core/drives/redis
```

## Broker

```go
package main

import (
    "github.com/gojuukaze/YTask/v3/drives/redis"
)

func main() {
	broker := redis.NewRedisBroker("127.0.0.1", "6379", "", 0, 3)
	// ...
}
```


## Backend

```go
package main

import (
    "github.com/gojuukaze/YTask/v3/drives/redis"
)

func main() {
	backend := redis.NewRedisBackend("127.0.0.1", "6379", "", 0, 10)
	// ...
}
```
