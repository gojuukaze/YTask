package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"math/rand"
	"sync"
	"time"
)

var (
	AMQPNil          = errors.New("rabbitMq get nil")
	ErrNoIdleChannel = errors.New("rabbitMq no idle channel")
	GetChanTimeout   = 60 * time.Second
	GetChanSleepTime = 100 * time.Millisecond
)

type Client struct {
	uri          string
	declareQueue map[string]struct{}
	lock         sync.Mutex
	Conn         *amqp.Connection
	IdleChan     map[string]*amqp.Channel
	NumOpen      int
	MaxChannel   int
}

func NewRabbitMqClient(host, port, user, password, vhost string, maxChannel int) *Client {

	c := Client{
		uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost),
		declareQueue: make(map[string]struct{}),
		IdleChan:     make(map[string]*amqp.Channel),
		MaxChannel:   maxChannel,
	}
	var err error
	c.Conn, err = amqp.Dial(c.uri)
	if err != nil {
		panic("YTask: connect rabbitMq error : " + err.Error())
	}
	return &c

}

func (c *Client) GetChannel() (*amqp.Channel, error) {

	timeout := time.Now().Add(GetChanTimeout)
	for {
		channel, err := c.getChannel()
		if err == nil {
			return channel, nil
		} else if err != ErrNoIdleChannel {
			return nil, err
		}
		if time.Now().After(timeout) {
			return nil, err
		}
		time.Sleep(GetChanSleepTime)
	}

}

func (c *Client) getChannel() (*amqp.Channel, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for id, channel := range c.IdleChan {
		isBad := channel.IsClosed()
		delete(c.IdleChan, id)
		if isBad {
			c.closeChan(id, channel)
		} else {
			return channel, nil
		}
	}
	if c.NumOpen == c.MaxChannel {
		return nil, ErrNoIdleChannel
	}
	// 创建新channel
	c.NumOpen++
	if c.Conn.IsClosed() {
		var err error
		c.Conn, err = amqp.Dial(c.uri)
		if err != nil {
			return nil, err
		}
	}
	channel, err := c.Conn.Channel()
	return channel, err
}

func (c *Client) closeChan(id string, channel *amqp.Channel) {
	if id != "" {
		delete(c.IdleChan, id)
	}
	c.NumOpen--
	channel.Close()
}

func (c *Client) PutChannel(channel *amqp.Channel, isBad bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if isBad {
		c.closeChan("", channel)
		return
	}
	// 这个id其实没啥用，只是为了配合map使用，所以每次随机生成
	id := fmt.Sprintf("%v-%v", rand.Intn(10), time.Now().UnixMilli())
	c.IdleChan[id] = channel

}

// 创建队列，如果队列已经存在，则忽略
func (c *Client) QueueDeclare(queueName string, channel *amqp.Channel) error {

	_, ok := c.declareQueue[queueName]
	if ok {
		return nil
	}
	if c.Conn.IsClosed() {
		var err error
		c.Conn, err = amqp.Dial(c.uri)
		if err != nil {
			return err
		}
	}

	_, err := channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{"x-max-priority": 9})
	return err
}

func (c *Client) Get(queueName string) ([]byte, error) {
	channel, err := c.GetChannel()
	// 因为获取channel时会判断一次isBad，因此放回时不再管isBad
	defer c.PutChannel(channel, false)
	if err != nil {
		return nil, err
	}

	if err = c.QueueDeclare(queueName, channel); err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	for {
		msg, ok, err := channel.Get(queueName, true)
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
	channel, err := c.GetChannel()
	defer c.PutChannel(channel, false)
	if err != nil {
		return err
	}
	if err = c.QueueDeclare(queueName, channel); err != nil {
		return err
	}
	err = channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        value,
		Priority:    Priority,
	})
	return err
}
