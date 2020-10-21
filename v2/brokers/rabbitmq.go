package brokers

import (
	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
)

type RabbitMqBroker struct {
	client   *drive.RabbitMqClient
	host     string
	port     string
	user     string
	password string
	vhost    string
	//poolSize int
}

func NewRabbitMqBroker(host, port, user, password, vhost string) RabbitMqBroker {
	return RabbitMqBroker{
		host:     host,
		port:     port,
		password: password,
		user:     user,
		vhost:    vhost,
		//poolSize: 0,
	}
}

func (r *RabbitMqBroker) Activate() {
	client := drive.NewRabbitMqClient(r.host, r.port, r.user, r.password, r.vhost)
	r.client = &client
}

func (r *RabbitMqBroker) SetPoolSize(n int) {
	//r.poolSize = n
}
func (r *RabbitMqBroker) GetPoolSize() int {
	//return r.poolSize
	return 0
}

func (r *RabbitMqBroker) Next(queueName string) (message.Message, error) {
	var msg message.Message
	value, err := r.client.Get(queueName)
	if err != nil {
		return msg, err
	}
	err = yjson.YJson.UnmarshalFromString(value, &msg)
	return msg, err
}

func (r *RabbitMqBroker) Send(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Publish(queueName, b, 0)
	return err
}

func (r *RabbitMqBroker) LSend(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Publish(queueName, b, 5)
	return err
}

func (r RabbitMqBroker) Clone() BrokerInterface {

	return &RabbitMqBroker{
		host:     r.host,
		port:     r.port,
		password: r.password,
		user:     r.user,
		//poolSize: 0,
	}
}
