工作流
===========

**注意！！使用工作流必须配置backend！！！** ，另外若backend出错会导致工作流终止。

开启工作流
---------------

通过 ``Workflow()`` 开启工作流，然后通过 ``Send()`` 提交子任务。

使用 ``Send()`` 时，只有第一个任务才需要任务参数，后续任务的参数默认为前一个任务的返回值。

.. code:: go

   ser := ytask.Server.NewServer(...)
   client:= ser.GetClient()


   tId, _ := client.Workflow().
   		Send("group1", "add", 123, 44).
   		Send("group1", "add").
   		Done()

工作流同样支持通过 ``SetTaskCtl()`` 设置任务参数。（ 设置延时任务时，只支持 ``RunAfter`` ）

对于下面的样例，第二个任务只有ExpireTime，不会继承第一个任务的设置

.. code:: go

   tId, _ := client.Workflow().
      		SetTaskCtl(client.RunAfter, 2*time.Second).
      		Send("group1", "add", 123, 44).
      		SetTaskCtl(client.ExpireTime, time.Now()).
      		Send("group1", "add").
      		Done()


获取工作流结果
-----------------

你可以使用 ``GetResult()`` 获取最后一个任务的返回值。

也可以使用  ``GetResult2()`` 获取任务流的运行进度。

.. code:: go

	result, err := client.GetResult2(id, time.Second*2, time.Millisecond*300)

	// 说明任务还未开始
	if yerrors.IsEqual(err, yerrors.ErrTypeTimeOut) {
		// ...
	}

	// status: waiting , running , success , failure , expired , abort
	for name, status:= range result.Workflow{
		fmt.Println(name, status)
	}
