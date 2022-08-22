
```shell
docker run --name ytask-memcache -d -p 11211:11211 memcached
cd drives/memcache
go test -v test/*.go

docker stop ytask-memcache
docker rm ytask-memcache

```