# RabbitMQ

## Installation

```shell
go get -u github.com/gojuukaze/YTask/core/drives/rabbitmq
```

## Broker

```go
package main

import (
    "github.com/gojuukaze/YTask/v3/drives/rabbitmq"
)

func main() {
	broker := rabbitmq.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest", "", 2)
	// ...
}
```

