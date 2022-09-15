中止任务
===========

你可以中止任务以及工作流，**注意！！使用中止任务必须配置backend！！！**

中止未运行任务
------------------------

通过 ``AbortTask()`` 设置中止任务标记。

第二个参数是中止标记的过期时间，若任务太多建议设置长一点，或设为-1后续手动清理。
（对于mongo这里设置过期时间是无效的，只能在NewBackend时设置）

.. code:: go

   client:= ser.GetClient()

   client.AbortTask(id, 10)

中止运行中的任务
------------------------

对于运行中的任务，你同样需要调用 ``AbortTask()`` 设置中止标记。然后在任务函数中手动检测，并中止。

.. code:: go

    func sendSMS(ctl *server.TaskCtl, userId, msg string) {
         
	phone := getUserPhone(userId)

    	if f, _ := ctl.IsAbort(); f {
    		ctl.Abort("手动中止")
    		// 别忘了return，否则会继续执行下去
    		return
    	}

    	Send(phone, msg)
    }
    
``IsAbort()`` 返回两个参数：``isAbort, err`` ，如果从backend获取数据出错 isAbort会设为false ，此时你要根据需求判断是否再次检测。
