延时任务
============
YTask并不能保证任务准时执行，
当你发现大量的延时任务到执行时间却没执行时，可尝试调大本地队列（v2.4+支持），调大YTask的任务并发数，或者部署更多的server端


开启延时任务
--------------

目前有两种方法开启延时任务：

1. 运行server时，把enableDelayServer设为true

.. code:: go

   ser := ytask.Server.NewServer(...)

   ser.Run("group", 10, true)

2. 通过config配置 (v2.4+)

.. code:: go

   config:=config.NewConfig(
       config.EnableDelayServer(true),
       config.delayServerQueueSize(50), // 本地队列大小
  )

提交延时任务
--------------

* RunAfter

.. code:: go

   client.SetTaskCtl(client.RunAfter, 1*time.Second).Send("group2", "add_sub", 123, 44)

* RunAt

.. code:: go

   runTime := time.Now().Add(1 * time.Second)
   client.SetTaskCtl(client.RunAt, runTime).Send("group2", "add_sub", 123, 44)

延时任务执行流程
------------------

首选说明两个概念：``inlineServer`` ：运行任务的server； ``delayServer`` : 获取延时任务的server。

1. 提交延时任务时YTask会在把任务提交到一个延时任务队列中（下面成为远程队列）。

2. delayServer会从远程队列中获取任务到本地队列并排序，
插入本地队列时，若本地队列已满，则会把离执行时间最远的任务出队，并重新插入远程队列的队尾。

3. delayServer定时取出本地队列中到执行时间的任务，并提交给inlineServer执行。也就是说延时任务与非延时任务公用一组并发的worker。

4. 但服务关闭时delayServer会把本地队列中的任务插入远程队列中（默认插到队头，如果broker不支持会插到队尾）
