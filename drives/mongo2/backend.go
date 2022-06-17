package mongo2

import (
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Backend struct {
	client     Client
	host       string
	port       string
	user       string
	password   string
	db         string
	collection string
	exTime     int
}

// NewMongoBackend
//  - exTime: Expiration time in seconds.  <=0: no expiration.
func NewMongoBackend(host, port, user, password, db, collection string, exTime int) Backend {
	return Backend{
		host:       host,
		port:       port,
		user:       user,
		password:   password,
		db:         db,
		collection: collection,
		exTime:     exTime,
	}
}

func (r *Backend) Activate() {
	r.client = NewMongoClient(r.host, r.port, r.user, r.password, r.db, r.collection, r.exTime)

}

func (r *Backend) SetPoolSize(n int) {
}

func (r *Backend) GetPoolSize() int {
	return 0
}

func (r *Backend) SetResult(result message.Result, exTime int) error {

	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}

	err = r.client.Set(result.GetBackendKey(), b)

	return err
}

func (r *Backend) GetResult(key string) (message.Result, error) {
	var result message.Result

	res, err := r.client.Get(key)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal(res.Data, &result)

	return result, err
}

func (r Backend) Clone() backends.BackendInterface {
	return &Backend{
		host:       r.host,
		port:       r.port,
		password:   r.password,
		db:         r.db,
		collection: r.collection,
		exTime:     r.exTime,
	}
}
