# Redis 

## Installation

```shell
go get -u github.com/gojuukaze/YTask/drives/redis/v3
```

## Broker

```go
package main

import (
    "github.com/gojuukaze/YTask/drives/redis/v3"
)

func main() {
	broker := redis.NewRedisBroker([]string{"127.0.0.1:6379"}, "", 0, 3, 0)
	// ...
}

```


## Backend

```go
package main

func main() {
	backend := redis.NewRedisBackend([]string{"127.0.0.1"}, "", 0, 10, 0)
	// ...
}

```
