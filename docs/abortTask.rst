中止任务
===========


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

对于运行中的任务，需要在任务函数中手动检测，并中止。

.. code:: go

	func abortWorker(ctl *server.TaskCtl, a int) int {
        // do ...

    	if f, _ := ctl.IsAbort(); f {
    		ctl.Abort("手动中止")
    		// 别忘了return
    		return 0
    	}

    	return a * a
    }
