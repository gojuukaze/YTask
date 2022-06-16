
```shell
docker run --name ytask-mongo -d -p 27017:27017 mongo
cd drives/mongo
go test -v test/*.go

docker stop ytask-mongo
docker rm ytask-mongo

```