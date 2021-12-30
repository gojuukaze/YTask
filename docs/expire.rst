任务过期时间
==============

目前只支持在client端设置过期时间，若任务发生重试且重试开始的时间超过了过期时间，任务会直接终止

.. code:: go

   taskId, err :=client.SetTaskCtl(client.ExpireTime,time.Now().Add(4*time.Second)).Send("group1","add",12,33)

   result, _ := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)

   if result.Status == message.ResultStatus.Expired{
         // do ...
   }

