package brokers

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"time"
)

type RocketMqBroker struct {
	client   *drive.RocketMqClient
	host     string
	port     string
}

func NewRocketMqBroker(host, port string) RocketMqBroker {
	 /*
	    1、目前不能自动创建topic
	    2、rocketmq topic名称不允许存在 ‘:’ ,
	    所以在生产、消费前先做了名称转换topic RocketMqClient.topicChecker 将非法字符全部转换为 ‘_’
		3、为提供pullConsumer实现，所以添加了在worker和consumer之间添加了 RocketMqClient.MsgChan
	    4、consumerOffset不能同步更新，所以任务执行时间更长
	    5、未支持RocketMqBroker.LSend
	 */
	return RocketMqBroker{
		host:     host,
		port:     port,
		//poolSize: 0,
	}
}
func (r *RocketMqBroker) Activate() {

	client := drive.NewRocketMqClient(r.host, r.port)
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

	queue,err:=r.client.Register(topic)
	if err!=nil{
		return msg, err
	}
	select {
	case value=<-queue:

	case <-time.After(2 * time.Second):
		return msg,yerrors.ErrEmptyQuery{}
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
		host:     r.host,
		port:     r.port,

		//poolSize: 0,
	}
}
func (r RocketMqBroker)Shutdown(){
	for topic,producer:=range r.client.RocketMqProducerMap{
		TRY1:
		err:=producer.Shutdown()
		if err !=nil{
			fmt.Println(topic,"producer shutdown err",err)
			goto TRY1
		}
	}
	for topic,consumer:=range r.client.RocketMqConsumerMap{
		close(r.client.MsgChanMap[topic])
		TRY2:
		err:=consumer.Shutdown()
		if err !=nil{
			fmt.Println(topic,"consumer shutdown err",err)
			goto TRY2
		}
	}
}