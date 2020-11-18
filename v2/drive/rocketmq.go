package drive

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"regexp"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type RocketMqClient struct {
	addr primitive.NamesrvAddr
	group string
	RocketMqProducerMap map[string]rocketmq.Producer
	RocketMqConsumerMap map[string]rocketmq.PushConsumer
	MsgChanMap map[string]chan string
}

func NewRocketMqClient(host,port string) RocketMqClient{

	var err error
	addr,err:=primitive.NewNamesrvAddr(host+":"+port)
	if err!=nil {
		panic("YTask: rocketMq error : " + err.Error())
	}
	return RocketMqClient{
		addr:addr,
		RocketMqProducerMap: make(map[string]rocketmq.Producer) ,
		RocketMqConsumerMap: make(map[string]rocketmq.PushConsumer),
		MsgChanMap: make(map[string]chan string),
	}
}


func (c *RocketMqClient) topicChecker(topic string)(string) {
	//rocketmq topic 只能包含%数字大小写字母及下划线和中划线
	re := regexp.MustCompile("[^A-z0-9_-]")
	//所以用下划线替换非法字符
	return re.ReplaceAllString(topic, "_")
}
func (c *RocketMqClient) Register(topic string) (<-chan string,error){
	topic=c.topicChecker(topic)


	if _,ok:=c.MsgChanMap[topic];!ok{
		c.MsgChanMap[topic]=make(chan string,0)
		output,err:=rocketmq.NewPushConsumer(
			consumer.WithNameServer(c.addr),
			consumer.WithGroupName(topic),
		)
		c.RocketMqConsumerMap[topic]=output
		/*addr,_:=internal.NewNamesrv(c.addr)
		options:=internal.ClientOptions{
			GroupName: topic,
			Namesrv: addr,
		}
		callBackChan:=make(chan interface{})
		client:=internal.GetOrNewRocketMQClient(options,callBackChan)
		offset:=consumer.NewRemoteOffsetStore(topic,client,addr)*/
		output.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			for _,msg := range msgs {
				fmt.Println("consumer:",string(msg.Body))
				c.MsgChanMap[topic]<-string(msg.Body)
			}
			return consumer.ConsumeSuccess, nil
		})
		err=output.Start()
		if err!=nil{
			fmt.Println("consumer start error ",err.Error())
		}
		return c.MsgChanMap[topic],err
	}
	//pull方式未实现
	//ref:=reflect.ValueOf(c.rocketMqConsumer)
	//method:=ref.MethodByName("Pull")
	//args:=[]reflect.Value{reflect.ValueOf(context.Background()),
	//	reflect.ValueOf(topic),
	//	reflect.ValueOf(consumer.MessageSelector{}),
	//	reflect.ValueOf(1)}
	//result:=method.Call(args)
	//res,err:=result[0].Interface().((*primitive.PullResult)),result[1].Interface().(error)
	return c.MsgChanMap[topic],nil
}
func (c *RocketMqClient) Publish(topic string,value interface{}, Priority uint8) error {

	topic=c.topicChecker(topic)
	if _,ok:=c.RocketMqProducerMap[topic];!ok{
		input, err := rocketmq.NewProducer(
			producer.WithNameServer(c.addr),
			producer.WithCreateTopicKey(topic),
			producer.WithGroupName(topic),
		)
		err=input.Start()
		if err!=nil {
			panic("YTask: rocketMq error : " + err.Error())
			return err
		}
		c.RocketMqProducerMap[topic]=input
	}
	fmt.Println("product:",string(value.([]byte)))
	_, err := c.RocketMqProducerMap[topic].SendSync(context.Background(),
		primitive.NewMessage(topic,value.([]byte)))
	if err != nil {
		return err
	}
	return nil
}
