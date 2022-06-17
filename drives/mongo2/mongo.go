package mongo2

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Result struct {
	Id         string    `json:"_id" bson:"_id"`
	Data       []byte    `json:"data" bson:"data"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

type Client struct {
	Uri        string
	DB         string
	Collection string
	Expires    int
}

func NewMongoClient(host, port, user, password, db, collection string, expires int) Client {
	var uri string
	if user != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
	}
	client := Client{uri, db, collection, expires}
	err := client.Init()
	if err != nil {
		panic("YTask: init mongo error : " + err.Error())
	}
	return client
}

// =======================
// high api
// =======================
func (c *Client) Get(key string) (Result, error) {
	var result Result
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := c.GetClient(ctx)

	if err != nil {
		return result, err
	}
	defer client.Disconnect(ctx)
	col := c.GetCollection(client)
	err = col.FindOne(ctx, bson.D{{"_id", key}}).Decode(&result)
	// 由于mongo不是立即清理过期数据，所以这里需要判断是否过期
	if err == nil && c.Expires > 0 && result.CreateTime.Add(time.Duration(c.Expires)*time.Second).Before(time.Now()) {
		err = mongo.ErrNoDocuments
	}

	return result, err

}

func (c *Client) Set(key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := c.GetClient(ctx)

	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	col := c.GetCollection(client)
	// 这里有个问题，如果doc在find之后过期，则ReplaceOne会报错。不过一般不用担心，key没那么容易重复

	filter := bson.D{{"_id", key}}
	err = col.FindOne(ctx, filter).Err()
	if err == mongo.ErrNoDocuments {
		_, err = col.InsertOne(ctx, Result{key, value, time.Now()})
		return err

	} else if err == nil {
		_, err = col.ReplaceOne(ctx, filter, Result{key, value, time.Now()})
		return err
	}

	return err
}

func (c *Client) GetClient(ctx context.Context) (*mongo.Client, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Uri))
	return client, err
}

func (c *Client) GetCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(c.DB).Collection(c.Collection)
}

func (c *Client) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := c.GetClient(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	col := c.GetCollection(client)

	if c.Expires > 0 {
		err = c.InitIndex(ctx, col.Indexes())
		if err != nil {
			return err
		}
	}

	return err
}

func (c *Client) InitIndex(ctx context.Context, index mongo.IndexView) error {

	cur, err := index.List(ctx)
	if err != nil {
		return err
	}
	// 索引为空则创建
	if !cur.Next(ctx) {
		indexOps := options.Index()
		indexOps.SetExpireAfterSeconds(int32(c.Expires))

		_, err = index.CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{"create_time", 1}},
			Options: indexOps,
		})
		return err
	}
	return nil

}
