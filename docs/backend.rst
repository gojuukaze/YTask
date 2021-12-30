Backend
===========
Backend用于保存任务结果

redisBackend
--------------

.. code:: go

   import "github.com/gojuukaze/YTask/v2"

   // 127.0.0.1 : host
   // 6379 : port
   // "" : password
   // 0 : db
   // 10 : 连接池大小.
   //      对于server端，如果为0，则值为min(10, numWorkers)
   //      对于client端, 你需要根据情况自行设置连接池
   ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 10)

memCacheBackend
---------------------

.. code:: go

   import "github.com/gojuukaze/YTask/v2"

   // 127.0.0.1 : host
   // 11211 : port
   // 10 : 连接池大小.
   //      对于server端，如果为0，则值为min(10, numWorkers)
   //      对于client端, 你需要根据情况自行设置连接池

   ytask.Backend.NewMemCacheBackend("127.0.0.1", "11211", 10)

mongoBackend
--------------

不支持设置过期时间，0代表不存储，>0代表永久存储

.. code:: go

   import "github.com/gojuukaze/YTask/v2"

   // 127.0.0.1 : host
   // 27017 : port
   // "" : username
   // "" : password
   // "task": db
   // "taks": collection

   ytask.Backend.NewMongoBackend("127.0.0.1", "27017", "", "", "task", "task")

自定义backend
----------------

你可以自行定义backend。同broker一样，调用\ ``Activate()``\ 时再建立连接。

.. code:: go

   type BackendInterface interface {
       SetResult(result message.Result, exTime int) error
       GetResult(key string) (message.Result, error)
       // Activate connection
       Activate()
       SetPoolSize(int)
       GetPoolSize() int
       Clone() BackendInterface

   }
