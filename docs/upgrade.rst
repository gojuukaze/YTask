升级说明
=============

从v2升级到v3
-----------------

你需要注意以下改动

* v3版本需要独立安装broker, backend

* TaskCtl结构体的位置迁移到了server包中，若需要获取RetryCount，则使用 ``GetRetryCount()``

  .. code:: go

     import "github.com/gojuukaze/YTask/v3/server"

     func add(ctl *server.TaskCtl, a, b int) int {

     	if ctl.GetRetryCount()==1 {
     		return 0
     	}
     	return a +b
     }

* v3版本修改了队列名中之前拼错的单词，及msg结构体。升级v3后你可以保持之前的server运行一段时间再关闭。
  或者使用 ``UseV2Name()`` 兼容

  .. code:: go

     import "github.com/gojuukaze/YTask/v3"

     ytask.UseV2Name()

     ser:=ytask.NewServer(...)

* 若你之前使用MongoBackend又不想丢失之前保存的任务结果，则使用 ``mongo`` 包，否则建议使用 ``mongo2`` 包

* v3移除了RocketMq支持，具体说明见：https://github.com/gojuukaze/YTask/tree/master/drives/rocketmq
