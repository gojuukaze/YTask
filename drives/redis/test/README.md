
```shell
docker run --name ytask-redis -d -p 6379:6379 redis
cd drives/redis
go test -v test/*.go

docker stop ytask-redis
docker rm ytask-redis

```

```shell
cd drives/redis
go test -v test/redisBackend_cluster_test.go
go test -v test/redisBroker_cluster_test.go
```