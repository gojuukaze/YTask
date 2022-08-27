客户端
=========

获取连接
----------

获取连接前一样需要初始化Server，然后调用\ ``GetClient()``\ 。\ ``NewServer``\ 的参数可以和服务端不同，但建议使用相同的参数配置

.. code:: go

   import "github.com/gojuukaze/YTask/v3"

   ser := ytask.Server.NewServer(
           ytask.Config.Broker(&broker),
           ytask.Config.Backend(&backend),
           ...
   )

   client = ser.GetClient()

发送信息
----------

| 使用\ ``Send``\ 发送任务信息，函数前两个参数为组名、任务，后面的参数是任务函数的参数。函数第一个返回值为任务id，可以用来获取任务结果。
| 发送消息时可以使用\ ``SetTaskCtl()``\ 配置该次任务的重试次数等

.. code:: go

   // group1 : 组名
   // add : 任务名
   // 12,33 ... : 任务参数
   // return :
   //   - taskId : taskId
   //   - err : error
   taskId,err:=client.Send("group1","add",12,33)

   // set retry count
   taskId,err=client.SetTaskCtl(client.RetryCount, 5).Send("group1","add",12,33)

   // set delay time
   taskId,err=client.SetTaskCtl(client.RunAfter, 2*time.Second).Send("group1","add",12,33)

   // set expire time
   taskId,err=client.SetTaskCtl(client.ExpireTime,time.Now().Add(4*time.Second)).Send("group1","add",12,33)

获取结果
----------
可通过 ``GetResult()`` ， ``GetResult2()``  获取结果

* ``GetResult()`` : 只有任务结束才返回（任务失败、完成都是结束）
* ``GetResult2()`` : backend中有结果就返回（一般来说任务开始执行就好有），这个通常用于获取任务流进度

---

| \ ``GetResult()``\ ， \ ``GetResult2()``\ 的第2个参数为超时时间，第3个参数为重新获取时间。
| 获取结果后可调用\ ``GetXX()``\ ，\ ``Get()``\ ，\ ``Gets()``\ 获取任务函数的返回结果。

.. code:: go

   // taskId :
   // 3*time.Second : timeout
   // 300*time.Millisecond : sleep time
   result, _ := client.GetResult(taskId, 3*time.Second, 300*time.Millisecond)

   if result.IsSuccess(){
       // get worker func return
       a,err:=result.GetInt64(0)
       b,err:=result.GetBool(1)

       // or
       var a int
       var b bool
       err:=result.Get(0, &a)
       err:=result.Get(1, &b)

       // or
       var a int
       var b bool
       err:=result.Gets(&a, &b)
   }



..

   | **重要！！**
   | YTask虽然提供获取结果功能，但不要过渡依赖。
   | 如果backend出错导致无法保存结果，YTask不会再次重试。因为对任务状态、结果的保存与运行任务的goroutine是同一个，不断重试会导致worker被占用。
     YTask优先保障任务运行，而不是结果保存。
   | 如果你特别需要任务结果，推荐你在任务函数中自行保存。
