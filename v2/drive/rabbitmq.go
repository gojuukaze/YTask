package drive

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)



type amqpErr string

func (e amqpErr) Error() string { return string(e) }

const AMQPNil = amqpErr("amqp: nil")

type RabbitMqClient struct {
	rabbitMqConn *amqp.Connection
	rabbitMqChan *amqp.Channel
	queueName    map[string]struct{}
}

func NewRabbitMqClient(host, port, user, password, vhost string) RabbitMqClient {
	client, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost))
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
		queueName:    make(map[string]struct{}),
	}

}

// =======================
// high api
// =======================

func (c *RabbitMqClient) queueDeclare(queueName string) error {
	_, ok := c.queueName[queueName]
	if ok {
		return nil
	}
	_, err := c.rabbitMqChan.QueueDeclare(queueName, true, false, false, false, amqp.Table{"x-max-priority": 9})
	return err
}

func (c *RabbitMqClient) Get(queueName string) (string, error) {
	if err := c.queueDeclare(queueName); err != nil {
		return "", err
	}
	msg, ok, err := c.rabbitMqChan.Get(queueName, true)
	if err!=nil{
		return "", err
	}
	if ok {
		return string(msg.Body), nil
	}else {
		return "", AMQPNil
	}
}

func (c *RabbitMqClient) Publish(queueName string, value interface{}, Priority uint8) error {
	if err := c.queueDeclare(queueName); err != nil {
		return err
	}
	err := c.rabbitMqChan.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        value.([]byte),
		Priority:    Priority,
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
