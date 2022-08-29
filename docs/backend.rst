Backend
===========
Backend用于保存任务结果

Redis
--------------
https://github.com/gojuukaze/YTask/tree/master/drives/redis

.. code:: go

   import "github.com/gojuukaze/YTask/drives/redis/v3"

   // 127.0.0.1 : host
   // 6379 : port
   // "" : password
   // 0 : db
   // 10 : 连接池大小.
   //      对于server端，如果为0，则值为min(10, numWorkers)
   //      对于client端, 你需要根据情况自行设置连接池
   redis.NewRedisBackend("127.0.0.1", "6379", "", 0, 10)

MemCache
---------------------
https://github.com/gojuukaze/YTask/tree/master/drives/memcache

.. code:: go

   import "github.com/gojuukaze/YTask/drives/memcache/v3"

   // []string{"127.0.0.1:11211"} : ["host:port"]
   // 10 : 连接池大小.
   //      对于server端，如果为0，则值为min(10, numWorkers)
   //      对于client端, 你需要根据情况自行设置连接池

   b := memcache.NewMemCacheBackend([]string{"127.0.0.1:11211"}, 10)

Mongo
--------------

https://github.com/gojuukaze/YTask/tree/master/drives/mongo2

如果你是从v2版升级，且希望保留旧数据，则使用 `mongo <https://github.com/gojuukaze/YTask/tree/master/drives/mongo>`__  ，否则使用mongo2

---

**注意1：** MongoBackend 无需设置连接池。

**注意2：** 若需设置过期时间，只能在 ``NewMongoBackend()`` 函数中设置，且设置后后续无法修改。
要修改的话只能手动修改MongoDB对应表的索引（修改ExpireAfterSeconds）

.. code:: go

   import "github.com/gojuukaze/YTask/drives/mongo2/v3"

   // 127.0.0.1 : host
   // 27017 : port
   // "" : username
   // "" : password
   // "test": db
   // "test": collection
   // 20: 过期时间，单位秒

	backend := mongo2.NewMongoBackend("127.0.0.1", "27017", "", "", "test", "test", 20)

自定义backend
----------------

你可以自定义backend。同broker一样，调用\ ``Activate()``\ 时再建立连接。 具体注意事项参照 :ref:`自定义Broker<custom>`

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
