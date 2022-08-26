package rabbitmq

import (
	"github.com/gojuukaze/YTask/v3/core/brokers"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/util/yjson"
	"github.com/gojuukaze/YTask/v3/core/yerrors"
)

type Broker struct {
	client   *Client
	host     string
	port     string
	user     string
	password string
	vhost    string
	poolSize int
}

func NewRabbitMqBroker(host, port, user, password, vhost string, poolSize int) Broker {
	return Broker{
		host:     host,
		port:     port,
		password: password,
		user:     user,
		vhost:    vhost,
		poolSize: poolSize,
	}
}

func (r *Broker) Activate() {
	r.client = NewRabbitMqClient(r.host, r.port, r.user, r.password, r.vhost, r.poolSize)
}

func (r *Broker) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *Broker) GetPoolSize() int {
	return r.poolSize
}

func (r *Broker) Next(queueName string) (message.Message, error) {
	var msg message.Message

	b, err := r.client.Get(queueName)
	if err != nil {
		if err == AMQPNil {
			err = yerrors.ErrEmptyQueue{}
		}
		return msg, err
	}

	err = yjson.YJson.Unmarshal(b, &msg)
	return msg, err
}

func (r *Broker) Send(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Publish(queueName, b, 0)
	return err
}

func (r *Broker) LSend(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Publish(queueName, b, 5)
	return err
}

func (r Broker) Clone() brokers.BrokerInterface {

	return &Broker{
		host:     r.host,
		port:     r.port,
		password: r.password,
		user:     r.user,
		vhost:    r.vhost,
		poolSize: r.poolSize,
	}
}
