Broker
==========

| YTask使用broker与任务队列通信，发送或接收任务。
| 支持的broker有：

redisBroker
--------------

.. code:: go

   import "github.com/gojuukaze/YTask/v2"

   // 127.0.0.1 : host
   // 6379 : port
   // "" : password
   // 0 : db
   // 10 : client连接池大小. (server端无需设置)
   //      对于client端, 你需要根据情况自行设置连接池
   ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 10)

rabbitMqBroker
-----------------

.. code:: go

   import "github.com/gojuukaze/YTask/v2"
   // 127.0.0.1 : host
   // 5672 : port
   // guest : username
   // guest : password

   ytask.Broker.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest", "")

rocketMqBroker
-----------------

不建议在延时任务中使用

.. code:: go

   import "github.com/gojuukaze/YTask/v2"

   ytask.Broker.NewRocketMqBroker([]string{"127.0.0.1:9876"},[]string{"127.0.0.1:10911"})

自定义broker
--------------

你可以自行定义broker。需要注意，因为系统中会调用\ ``SetPoolSize``\ 设置连接池，所以初始化broker时不要建立连接，调用\ ``Activate()``\ 时再建立。

如果你的broker不支持连接池，那可以不用管Activate,SetPoolSize,GetPoolSize三个方法，直接返回空就行。

获取任务时，应不断循环获取，而不是阻塞在这，若队列为空，则返回 ``ErrTypeEmptyQuery`` 错误。

.. code:: go

   type BrokerInterface interface {
       // 获取任务
       Next(queryName string) (message.Message, error)
       // 发送任务
       Send(queryName string, msg message.Message) error
       // 把任务插到队头
       //  - 如果你的broker不支持，也没有优先队列这样的替代方案，则可以复用Send（这样做会影响延时任务的处理时间）。
       LSend(queryName string, msg message.Message) error
       // 建立连接
       Activate()
       SetPoolSize(int)
       GetPoolSize()int
       // 用当前broker的配置生成个新的broker
       Clone() BrokerInterface
   }
