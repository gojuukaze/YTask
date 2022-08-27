Error
==========

内置的错误类型
-------------------

.. code:: go

   const (
	  ErrTypeEmptyQueue      = 1 // 队列为空， broker获取任务时用到
	  ErrTypeUnsupportedType = 2 // 不支持此参数类型
	  ErrTypeOutOfRange      = 3 // 暂时没用
	  ErrTypeNilResult       = 4 // 任务结果为空
	  ErrTypeTimeOut         = 5 // broker，backend超时
	  ErrTypeServerStop      = 6 // 服务已停止
	  ErrTypeSendMsg    = 7 // 通过broker发送消息失败，目前工作流中发送下一个任务时会用到
	  ErrTypeNilBackend = 8
	  ErrTypeAbortTask = 9
   )


比较错误
-------------

.. code:: go

   import  "github.com/gojuukaze/YTask/v3/yerrors"
   yerrors.IsEqual(err, yerrors.ErrTypeNilResult)
