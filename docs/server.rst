服务端
======

初始化
---------

-  使用\ ``NewServer()``\ 初始化服务，其参数是server的配置，所有配置在下面

.. code:: go

   import "github.com/gojuukaze/YTask/v2"

   ser := ytask.Server.NewServer(
           ytask.Config.Broker(&broker),
           ytask.Config.Backend(&backend),
           ...
   )

服务端配置
--------------
+--------+----+----+------------------+------------------------------+
| 配置   | 是 | 默 | code             | 说明                         |
|        | 否 | 认 |                  |                              |
|        | 必 | 值 |                  |                              |
|        | 须 |    |                  |                              |
+========+====+====+==================+==============================+
| Broker | \* |    | yta              |                              |
|        |    |    | sk.Config.Broker |                              |
+--------+----+----+------------------+------------------------------+
| B      |    | n  | ytas             |                              |
| ackend |    | il | k.Config.Backend |                              |
+--------+----+----+------------------+------------------------------+
| Debug  |    | F  | yt               | 是否开启debug                |
|        |    | AL | ask.Config.Debug |                              |
|        |    | SE |                  |                              |
+--------+----+----+------------------+------------------------------+
| S      |    | 1d | ytask.Conf       | 单                           |
| tatusE |    | ay | ig.StatusExpires | 位：秒，任务状态的过期时间,  |
| xpires |    |    |                  | -1:永久保                    |
|        |    |    |                  | 存（有的backend可能不支持）  |
+--------+----+----+------------------+------------------------------+
| R      |    | 1d | ytask.Conf       | 单                           |
| esultE |    | ay | ig.ResultExpires | 位：秒，任务结果的过期时间,  |
| xpires |    |    |                  | -1:永久保存                  |
|        |    |    |                  | （有的backend可能不支持）    |
+--------+----+----+------------------+------------------------------+

-  任务状态、结果有什么不同？

   -  状态： 任务的开始、运行、成功、失败状态
   -  结果： 函数的返回值

-  对于\ ``mongo backend`` 过期时间 0代表不存储，>0代表永久存储

注册任务
--------------

使用\ ``Add``\ 注册任务

.. code:: go

   func addFunc(a,b int) (int, bool){
       return a+b, true
   }

   // group1 : 任务所属组，也是队列的名字
   // add : 任务名
   // addFunc : 任务函数
   ser.Add("group1","add",addFunc)

任务函数的 参数、返回值 支持所有能被系列化为json的类型。

如果需要在函数中控制任务的重试等东西，则函数的第一个参数为\ ``TaskCtl``
(带\ ``TaskCtl``\ 函数和其他函数使用上没区别)

如：

.. code:: go

   func addFunc(ctl *controller.TaskCtl, a int, b int) (int, int) {
       if ... {
           // retry
           ctl.Retry(errors.New("ctl.Retry"))
           return 0, 0
       }

       return a + b, a - b
   }

   ser.Add("group1","add",addFunc)

任务回调
--------------
v2.4+支持

注册任务时可以为任务添加回调函数，回调函数在 **任务结束** 后调用，前几个参数为任务的参数，最后一个参数为返回结果。

回调函数报错并不会影响任务的结果，另外由于回调函数和任务函数是在同一个goroutine中执行，回调函数不结束会导致当前worker一直被占用，
因此你需要根据实际需求评估回调函数需要执行的任务。

.. code:: go

   func addFunc(a,b int) (int, bool){
       return a+b, true
   }

   func callbackFunc(a,b int, result *message.Result) {
       if result.IsSuccess(){
          // do ...
       }else {
         // do ...
       }
   }

   ser.Add("group1", "add", addFunc, callbackFunc)


运行与停止
--------------

.. code:: go

   // group1 : 运行的组名
   // 3 : 并发任务数
   // false : 是否开启延时任务
   ser.Run("group1", 3, false)

   quit := make(chan os.Signal, 1)
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   <-quit
   ser.Shutdown(context.Background())

运行多个group
--------------

-  V2.2.0+ 才支持

.. code:: go

   ser:=ytask.Server.NewServer(...)

   ser.Run("g1", 5)
   ser.Run("g2", 5)
