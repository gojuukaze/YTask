# MongoDB

**【注意】使用MongoDB时，若需设置过期时间，只能在 ``NewMongoBackend()`` 函数中设置**

## Installation

```shell
go get -u github.com/gojuukaze/YTask/core/drives/mongo2
```

## Backend

```go
package main

import (
    "github.com/gojuukaze/YTask/v3/drives/mongo2"
)

func main() {
	backend := mongo2.NewMongoBackend("127.0.0.1", "27017", "", "", "test", "test", 2, 20)
	// ...
}
```
