package drive

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"regexp"
	"sync"
)

type RocketMqClient struct {
	options *clientOptions
	Producer rocketmq.Producer
	ConsumerMap map[string]rocketmq.PushConsumer
 	MsgChanMap map[string]chan string
	Admin admin.Admin
	topicMap sync.Map

}
type clientOptions struct {
	NamesrvAddr primitive.NamesrvAddr
	AutoCreateTopic bool
	BrokerAddr []string
	auto bool
}
type ClientOption func(options *clientOptions)

func defaultAdminOptions() *clientOptions {
	return &clientOptions{}
}
func WithNameSrvAddr(addr []string) ClientOption{
	return func(opts *clientOptions) {
		opts.NamesrvAddr=addr
	}
}
func WithBrokerAddr(addr []string) ClientOption{
	return func(opts *clientOptions) {
		opts.BrokerAddr=addr
	}
}
func WithAutoCreateTopic(auto bool) ClientOption {
	return func(opts *clientOptions) {
		opts.auto=auto
	}
}


func NewRocketMqClient(opts... ClientOption) RocketMqClient{
	defaultOpts := defaultAdminOptions()
	for _, opt := range opts {
		opt(defaultOpts)
	}
	var adm admin.Admin
	var err error

	adm, err = admin.NewAdmin(admin.WithResolver(
		primitive.NewPassthroughResolver(defaultOpts.NamesrvAddr)))
	if err!=nil {
		panic("YTask: admin create error : "+err.Error())
	}

	input, err := rocketmq.NewProducer(
		producer.WithNameServer(defaultOpts.NamesrvAddr),
		//producer.WithCreateTopicKey(topic),
	)
	if err!=nil {
		panic("YTask[RockerMQ]: Producer create error : " + err.Error())
	}
	err=input.Start()
	if err!=nil {
		panic("YTask[RockerMQ]: Producer start error : " +err.Error())
	}
	return RocketMqClient{
		options: defaultOpts,
		Producer: input,
		ConsumerMap: make(map[string]rocketmq.PushConsumer),
		MsgChanMap: make(map[string]chan string),
		Admin: adm,
	}
}
func (c *RocketMqClient) topicCreator(topic string) {
	if c.options.BrokerAddr == nil {
		return
	}
	//create topic
	for _, addr := range c.options.BrokerAddr {
		err := c.Admin.CreateTopic(
			context.Background(),
			admin.WithTopicCreate(topic),
			admin.WithBrokerAddrCreate(addr),
			admin.WithReadQueueNums(1),
			admin.WithWriteQueueNums(1),
			admin.WithPerm(6),
		)
		if err != nil {
			fmt.Println("YTask[RocketMQ]: create topic error:", err.Error())
		}
	}
}
func (c *RocketMqClient) TopicDeleter(topic string) {
	//delete topic
	err:=c.Admin.DeleteTopic(
		context.Background(),
		admin.WithTopicDelete(topic),
	)
	if err != nil {
		fmt.Println("YTask[RocketMQ]: delete topic error:", err.Error())
	}
}

func (c *RocketMqClient) topicChecker(topic string)(string) {
	//rocketmq topic 只能包含%数字大小写字母及下划线和中划线
	re := regexp.MustCompile("[^A-z0-9_-]")
	//所以用下划线替换非法字符
	return re.ReplaceAllString(topic, "_")
}

func (c *RocketMqClient) Register(topic string)(<-chan string,error){
	topic=c.topicChecker(topic)
	if _,ok:=c.ConsumerMap[topic];!ok{
		if _,ok:=c.topicMap.LoadOrStore(topic,1);!ok {
			if c.options.auto {
				c.topicCreator(topic)
			}
		}
		c.MsgChanMap[topic]=make(chan string,0)
		output,err:=rocketmq.NewPushConsumer(
			consumer.WithNameServer(c.options.NamesrvAddr),
			consumer.WithGroupName(topic),
		)
		if err!=nil {
			panic("YTask[RockerMQ]: Consumer create error : " + err.Error())
		}
		output.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			for _,msg := range msgs {
				fmt.Println("consumer:",string(msg.Body))
				c.MsgChanMap[topic]<-string(msg.Body)
			}
			return consumer.ConsumeSuccess, nil
		})
		err=output.Start()
		if err!=nil {
			panic("YTask[RockerMQ]: Consumer start error : " +err.Error())
		}
		c.ConsumerMap[topic]=output
		return c.MsgChanMap[topic],nil
	}

	return c.MsgChanMap[topic],nil
}


func (c *RocketMqClient) Publish(topic string,value interface{}, Priority uint8) error {
	if _,ok:=c.topicMap.LoadOrStore(topic,1);!ok {
		if c.options.auto {
			c.topicCreator(topic)
		}
	}
	topic=c.topicChecker(topic)
	fmt.Println("produce:",string(value.([]byte)))
	_, err := c.Producer.SendSync(context.Background(),
		primitive.NewMessage(topic,value.([]byte)))
	if err != nil {
		return err
	}
	return nil
}
