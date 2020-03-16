package drive

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMqClient struct {
	rabbitMqConn *amqp.Connection
	rabbitMqChan *amqp.Channel

}

func NewRabbitMqClient(host, port, user, password string) RabbitMqClient {
	client, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port))
	if err != nil {
		panic("YTask: connect rabbitMq error : " + err.Error())
	}

	channel, err := client.Channel()
	if err != nil {
		panic("YTask: get rabbitMq channel error : " + err.Error())
	}

	return RabbitMqClient{
		rabbitMqConn: client,
		rabbitMqChan: channel,
	}

}

// =======================
// high api
// =======================
func (c *RabbitMqClient) Get(key string) (string, error) {
	msg, ok, err := c.rabbitMqChan.Get(key, true)
	if ok && err == nil {
		return string(msg.Body), nil
	}
	return "", err
}

func (c *RabbitMqClient) Set(key string, value interface{}) error {
	q, _ := c.rabbitMqChan.QueueDeclare(key, false, false, false, false, nil)
	err := c.rabbitMqChan.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType:     "text/plain",
		Body:            value.([]byte),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *RabbitMqClient) Ping() error {
	closed := c.rabbitMqConn.IsClosed()
	if closed {
		return errors.New("rabbitMq connection is closed")
	}
	return nil
}

func (c *RabbitMqClient) Close() {
	_ = c.rabbitMqChan.Close()
	_ = c.rabbitMqConn.Close()
}