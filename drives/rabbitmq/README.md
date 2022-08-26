# RabbitMQ

## Installation

```shell
go get -u github.com/gojuukaze/YTask/drives/rabbitmq/v3
```

## Broker

```go
package main

import (
    "github.com/gojuukaze/YTask/drives/rabbitmq/v3"
)

func main() {
	broker := rabbitmq.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest", "", 2)
	// ...
}
```

