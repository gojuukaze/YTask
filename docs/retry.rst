任务重试
===========

触发重试
----------

有两种方法可以触发2重试

-  使用 panic

.. code:: go


   func add(a, b int){
       panic("xx")
   }

-  使用 TaskCtl

.. code:: go


   func add(ctl *controller.TaskCtl,a, b int){
       ctl.Retry(errors.New("xx"))
       return
   }

设置重试次数
--------------

默认的重试次数是3次，目前只支持在client端设置

-  in client

.. code:: go

   client.SetTaskCtl(client.RetryCount, 5).Send("group1", "retry", 123, 44)

禁用重试
---------

-  在server端针对某个任务禁用

.. code:: go

   func add(ctl *controller.TaskCtl,a, b int){
       ctl.SetRetryCount(0)
       return
   }

-  在client端对此次任务禁用

.. code:: go

   client.SetTaskCtl(client.RetryCount, 0).Send("group1", "retry", 123, 44)
