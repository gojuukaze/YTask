# MongoDB

**【注意】使用MongoDB时，若需设置过期时间，只能在 ``NewMongoBackend()`` 函数中设置**

## Installation

```shell
go get -u github.com/gojuukaze/YTask/drives/mongo2/v3
```

## Backend

```go
package main

import (
    "github.com/gojuukaze/YTask/drives/mongo2/v3"
)

func main() {
	backend := mongo2.NewMongoBackend("127.0.0.1", "27017", "", "", "test", "test", 200)
	// ...
}
```
