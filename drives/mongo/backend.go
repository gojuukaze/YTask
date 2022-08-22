package mongo

import (
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Backend struct {
	client     MongoClient
	host       string
	port       string
	user       string
	password   string
	db         string
	collection string
	//poolSize int
}

func NewMongoBackend(host, port, user, password, db, collection string) Backend {
	return Backend{
		host:       host,
		port:       port,
		user:       user,
		password:   password,
		db:         db,
		collection: collection,
		//poolSize: 0,
	}
}

func (r *Backend) Activate() {
	client := NewMongoClient(r.host, r.port, r.user, r.password, r.db, r.collection)
	r.client = client
}

func (r *Backend) SetPoolSize(n int) {
	//r.poolSize = n
}

func (r *Backend) GetPoolSize() int {
	//return r.poolSize
	return 0
}

func (r *Backend) SetResult(result message.Result, exTime int) error {

	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}

	err = r.client.Set(result.GetBackendKey(), b, exTime)

	return err
}

func (r *Backend) GetResult(key string) (message.Result, error) {
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

func (r Backend) Clone() backends.BackendInterface {
	return &Backend{
		host:       r.host,
		port:       r.port,
		password:   r.password,
		db:         r.db,
		collection: r.collection,
	}
}
