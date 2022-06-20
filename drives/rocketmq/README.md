# 不再支持！！

`rocketmq-client` 的设计与 `ytask` 的设计是不兼容的。 具体是 `broker.Next()` 的设计。  
v2版的这部分代码其实是有问题的，因此v3不再支持。

## ytask

`ytask` 的 `broker.Next()` 设计为:每次获取一条消息并返回，这样可以让 `ytask` 主动控制获取消息的频率（即有woker空闲时才获取消息）。

当然这牺牲了处理速度，且如果没有连接池或者保持链接，每次调用 `broker.Next()` 都会建立一个连接。 但这样做能让整个逻辑变得简单。另外这点开销也是能接受的。

## rocketmq-client

`rocketmq-client` 是典型的 "生产-消费" 模型，通过 `Consumer.Subscribe()` 订阅消息，有新消息就立即消费。

这导致无法实现 `ytask` 的 `broker.Next()` 。

如果你读过v2版的代码，会发现v2其实用了chan。启动时先创建一个无缓冲的chan，消费者获取消息后提交到chan，每次调用 `broker.Next()` 时从chan中获取一条消息并返回。

这样做会导致服务关闭时要做而外的清理工作。且为了防止多个消费者同时有消息等待提交到chan，`ReadQueueNums`被设置为了1，这导致只能部署一个服务。

---

综上所述，不再支持`rocketmq`。
