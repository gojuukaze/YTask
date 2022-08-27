
```shell
docker run --name ytask-rabbit -d -p 5672:5672 rabbitmq
cd drives/rabbitmq
# rabbitmq启动有点慢，这里要等一下
go test -v test/*.go

docker stop ytask-rabbit
docker rm ytask-rabbit

```