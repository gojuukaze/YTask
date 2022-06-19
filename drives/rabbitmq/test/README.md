
```shell
docker run --name ytask-rabbit -d -p 5672:5672 rabbitmq
cd drives/rabbitmq
go test -v test/*.go

docker stop ytask-rabbit
docker rm ytask-rabbit

```