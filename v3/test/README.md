# test

## run test

```shell
go test -v test/*.go
```

## 编写说明

使用``LocalBroker``与``LocalBackend`` 进行测试时要注意：

* 只能执行一次 ``Active()`` ，因为Active操作会清空数据
* ``ser.GetClient()`` 应该在 ``ser.Run()`` 之前调用，同样也只能执行一次GetClient操作
* 由于LocalBroker执行Next, Send时会争抢锁，有时候测试不通过可能是因为一直没抢到锁超时了，可以把此项测试移到redis中，或修改超时时间
* GetResult时，timeout建议设长点，sleepTime建议为100毫秒或更短