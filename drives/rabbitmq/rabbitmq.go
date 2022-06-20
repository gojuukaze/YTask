package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var AMQPNil = errors.New("rabbitMq get nil")

type Client struct {
	uri          string
	declareQueue map[string]struct{}
}
type rabbitSession struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (s rabbitSession) close() {
	s.channel.Close()
	s.conn.Close()
}
func NewRabbitMqClient(host, port, user, password, vhost string) Client {

	c := Client{
		uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost),
		declareQueue: make(map[string]struct{}),
	}
	_, err := c.getSession(c.uri)
	if err != nil {
		panic("YTask: connect rabbitMq error : " + err.Error())
	}
	return c

}

func (c *Client) getSession(uri string) (*rabbitSession, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &rabbitSession{conn, channel}, nil
}

// 创建队列，如果队列已经存在，则忽略
func (c *Client) queueDeclare(queueName string, session *rabbitSession) error {

	_, ok := c.declareQueue[queueName]
	if ok {
		return nil
	}
	_, err := session.channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{"x-max-priority": 9})
	return err
}

func (c *Client) Get(queueName string) ([]byte, error) {
	session, err := c.getSession(c.uri)
	if err != nil {
		return nil, err
	}
	defer session.close()
	if err = c.queueDeclare(queueName, session); err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	for {
		msg, ok, err := session.channel.Get(queueName, true)
		if err != nil {
			return nil, err
		}
		if ok {
			return msg.Body, nil
		}

		select {
		case <-ctx.Done():
			return nil, AMQPNil
		case <-time.After(time.Millisecond * 300):
		}
	}

}

func (c *Client) Publish(queueName string, value []byte, Priority uint8) error {
	session, err := c.getSession(c.uri)
	if err != nil {
		return err
	}
	defer session.close()
	if err = c.queueDeclare(queueName, session); err != nil {
		return err
	}
	err = session.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        value,
		Priority:    Priority,
	})
	return err
}
