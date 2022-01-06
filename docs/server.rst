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

.. 
 注意，rst的表格比较特殊，每列的字符数必须一致，直接编辑很可能会有语法错误。 
 编辑时先把下列代码转为 mediaWiki或者复制文档里的html，然后在这个网站 https://tableconvert.com/ 可视化编辑，再复制 mediaWiki，然后再转为rst。（md与rst互转会有问题）
 格式转换可通过这个网站：https://pandoc.org/try/ ，可在这个网站验证rst的正确性：http://rst.ninjs.org/

+-------------+----------+-------------+-------------+-------------+
| 配置        | 是否必须 | 默认值      | code        | 说明        |
+=============+==========+=============+=============+=============+
| Broker      | \*       |             | ytask.Co    |             |
|             |          |             | nfig.Broker |             |
+-------------+----------+-------------+-------------+-------------+
| Backend     |          | nil         | ytask.Con   |             |
|             |          |             | fig.Backend |             |
+-------------+----------+-------------+-------------+-------------+
| Logger      |          | log.        | ytask.Co    | logger,     |
|             |          | YTaskLogger | nfig.Logger | v2.5+支持   |
+-------------+----------+-------------+-------------+-------------+
| Debug       |          | FALSE       | ytask.C     | 是          |
|             |          |             | onfig.Debug | 否开启debug |
+-------------+----------+-------------+-------------+-------------+
| St          |          | 1day        | ytas        | 单位        |
| atusExpires |          |             | k.Config.St | ：秒，任务  |
|             |          |             | atusExpires | 状态的过期  |
|             |          |             |             | 时间,-1:永  |
|             |          |             |             | 久保存（有  |
|             |          |             |             | 的backend可 |
|             |          |             |             | 能不支持）  |
+-------------+----------+-------------+-------------+-------------+
| Re          |          | 1day        | ytas        | 单位        |
| sultExpires |          |             | k.Config.Re | ：秒，任务  |
|             |          |             | sultExpires | 结果的过期  |
|             |          |             |             | 时间,-1:永  |
|             |          |             |             | 久保存（有  |
|             |          |             |             | 的backend可 |
|             |          |             |             | 能不支持）  |
+-------------+----------+-------------+-------------+-------------+



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

日志接口
--------------
v2.5+支持

只要实现 log.LoggerInterface 接口即可，默认已经实现一个基于 logrus 的 logger

.. code:: go

	logger := ytask.Logger.NewYTaskLogger()

	Server := ytask.Server.NewServer(
	    ...
		ytask.Config.Logger(logger),		// 可以不设置 logger
		...
	)


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
