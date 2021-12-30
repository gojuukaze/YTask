package backends

import (
	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBackend struct {
	client   *drive.MongoClient
	host     string
	port     string
	user string
	password string
	db string
	collection string
	//poolSize int
}

func NewMongoBackend(host, port , user, password, db, collection string) MongoBackend {
	return MongoBackend{
		host:     host,
		port:     port,
		user: user,
		password: password,
		db: db,
		collection: collection,
		//poolSize: 0,
	}
}

func (r *MongoBackend) Activate() {
	client := drive.NewMongoClient(r.host, r.port, r.user, r.password, r.db, r.collection)
	r.client = &client
}

func (r *MongoBackend) SetPoolSize(n int) {
	//r.poolSize = n
}

func (r *MongoBackend) GetPoolSize() int {
	//return r.poolSize
	return 0
}

func (r *MongoBackend) SetResult(result message.Result, exTime int) error {

	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}

	err = r.client.Set(result.GetBackendKey(), b, exTime)

	return err
}

func (r *MongoBackend) GetResult(key string) (message.Result, error) {
	var result message.Result

	b, err := r.client.Get(key)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal([]byte(b), &result)

	return result, err
}

func (r MongoBackend) Clone() BackendInterface{
	return  &MongoBackend{
		host:     r.host,
		port:     r.port,
		password: r.password,
		db:       r.db,
		collection: r.collection,
	}
}