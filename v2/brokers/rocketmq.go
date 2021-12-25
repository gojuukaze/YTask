package brokers

import (
	"fmt"

	"time"

	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
	"github.com/gojuukaze/YTask/v2/yerrors"
)

type RocketMqBroker struct {
	client      *drive.RocketMqClient
	namesrvAddr []string
	brokerAddr  []string
	auto        bool
}

func NewRocketMqBroker(namesrvAddr []string, brokerAddr ...[]string) RocketMqBroker {
	/*
	   FIX：1、目前不能自动创建topic (mqadmin手动创建，并设置读写队列数为1)
	   2、consumerOffset不能同步更新，所以任务执行时间更长
	   (需要将队列中多余的message消费掉才能消费到当前taskId对应的消息)
	   3、未支持RocketMqBroker.LSend
	*/
	var auto bool
	if len(brokerAddr) > 0 {
		auto = true
	}
	return RocketMqBroker{
		namesrvAddr: namesrvAddr,
		brokerAddr:  brokerAddr[0],
		auto:        auto,
	}
}
func (r *RocketMqBroker) Activate() {
	client := drive.NewRocketMqClient(
		drive.WithNameSrvAddr(r.namesrvAddr),
		drive.WithBrokerAddr(r.brokerAddr),
		drive.WithAutoCreateTopic(r.auto))
	r.client = &client

}

func (r *RocketMqBroker) SetPoolSize(n int) {
	//r.poolSize = n
}
func (r *RocketMqBroker) GetPoolSize() int {
	//return r.poolSize
	return 0
}

func (r *RocketMqBroker) Next(topic string) (message.Message, error) {
	var msg message.Message
	var value string
	var err error

	queue, err := r.client.Register(topic)
	if err != nil {
		return msg, err
	}
	select {
	case value = <-queue:

	case <-time.After(5 * time.Second):
		return msg, yerrors.ErrEmptyQuery{}
	}

	err = yjson.YJson.UnmarshalFromString(value, &msg)
	return msg, err
}

func (r *RocketMqBroker) Send(topic string, msg message.Message) error {

	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Publish(topic, b, 0)
	return err
}

func (r *RocketMqBroker) LSend(queueName string, msg message.Message) error {
	// 未实现
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.Publish(queueName, b, 5)
	return err
}

func (r RocketMqBroker) Clone() BrokerInterface {

	return &RocketMqBroker{
		namesrvAddr: r.namesrvAddr,
		brokerAddr:  r.brokerAddr,
		auto:        r.auto,
	}
}

//目前不做使用
func (r RocketMqBroker) Shutdown() {
TRY1:
	err := r.client.Producer.Shutdown()
	if err != nil {
		fmt.Println("YTask[RocketMQ]: producer shutdown err:", err)
		goto TRY1
	}
	for topic, consumer := range r.client.ConsumerMap {
		err := consumer.Unsubscribe(topic)
		if err != nil {
			fmt.Println("YTask[RocketMQ]: Unsubscribe err: ", err)
		}
	TRY2:
		err = consumer.Shutdown()
		if err != nil {
			fmt.Println(topic, "YTask[RocketMQ]: consumer shutdown err: ", err)
			goto TRY2
		}
		close(r.client.MsgChanMap[topic])

		r.client.TopicDeleter(topic)
		//consumer.Shutdown()方法没法及时同步，所以在异步任务结束后删除topic
		//重新开启任务时创建topic,代理点位和消费点位重置为0
		//不得已为之，待改善
	}
	r.client.Admin.Close()
}
