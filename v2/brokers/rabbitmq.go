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
	//poolSize int
}

func NewRabbitMqBroker(host, port, user, password string) RabbitMqBroker {
	return RabbitMqBroker{
		host:     host,
		port:     port,
		password: password,
		user:     user,
		//poolSize: 0,
	}
}

func (r *RabbitMqBroker) Activate() {
	client := drive.NewRabbitMqClient(r.host, r.port, r.user, r.password)
	r.client = &client
}

func (r *RabbitMqBroker) SetPoolSize(n int) {
	//r.poolSize = n
}
func (r *RabbitMqBroker) GetPoolSize() int {
	//return r.poolSize
	return 0
}

func (r *RabbitMqBroker) Next(queryName string) (message.Message, error) {
	var msg message.Message
	value, err := r.client.Get(queryName)
	if err != nil {
		return msg, err
	}

	err = yjson.YJson.UnmarshalFromString(value, &msg)
	return msg, err
}

func (r *RabbitMqBroker) Send(queryName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Set(queryName, b)
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
