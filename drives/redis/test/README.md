
```shell
docker run --name ytask-redis -d -p 6379:6379 redis
cd drives/redis
go test -v test/*.go

docker stop ytask-redis
docker rm ytask-redis

```