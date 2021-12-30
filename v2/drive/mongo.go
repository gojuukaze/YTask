package drive

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Result struct {
	Uuid string `json:"uuid" bson:"uuid"`
	Res string `json:"res" bson:"res"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

type MongoClient struct {
	mongoConn *mongo.Client
	mongoCollection *mongo.Collection
	hasSetExTime bool
}

func NewMongoClient(host, port, user, password, db, collection string) MongoClient {
	var clientOptions *options.ClientOptions
	if user != "" {
		clientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port))
	} else {
		clientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
	}

	ctx, _ := context.WithTimeout(context.Background(), 15 * time.Second)
	//client, err := mongo.Connect(context.TODO(), clientOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic("YTask: connect mongo error : " + err.Error())
	}

	//err = client.Ping(context.TODO(), nil)
	err = client.Ping(ctx, nil)
	if err != nil {
		panic("YTask: connect mongo error : " + err.Error())
	}

	coll := client.Database(db).Collection(collection)

	return MongoClient{
		mongoConn: client,
		mongoCollection: coll,
	}
}

// =======================
// high api
// =======================
func (c *MongoClient) Get(key string) (string, error) {
	var res Result
	filter := bson.D{{"uuid", key}}
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	//err := c.mongoCollection.FindOne(context.TODO(), filter).Decode(&res)
	err := c.mongoCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.Res, nil
}

func (c *MongoClient) Set(key string, value interface{}, exTime int) error {
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)

	if !c.hasSetExTime {
		_, _ = c.mongoCollection.Indexes().DropOne(ctx, "create_time_ttl_index")
		if exTime != -1 {
			option := options.Index()
			option.SetName("create_time_ttl_index")
			option.SetExpireAfterSeconds(int32(exTime))
			option.SetBackground(true)

			_, _ = c.mongoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bsonx.Doc{{"create_time", bsonx.Int32(1)}},
				Options: option,
			})
		}
		c.hasSetExTime = true
	}

	res := Result{key, string(value.([]byte)), time.Now()}
	_ , err := c.Get(key)
	if err == mongo.ErrNoDocuments {
		//_, err := c.mongoCollection.InsertOne(context.TODO(), res)
		_, err := c.mongoCollection.InsertOne(ctx, res)
		return err
	} else {
		filter := bson.D{{"uuid", key}}

		//update := bson.D{
		//	{"$set", bson.D{
		//		{"uuid", key},
		//		{"res", string(value.([]byte))},
		//	}},
		//}
		//
		//_, err := c.mongoCollection.UpdateOne(context.TODO(), filter, update)

		// or
		//_, err := c.mongoCollection.ReplaceOne(context.TODO(), filter, res)
		_, err := c.mongoCollection.ReplaceOne(ctx, filter, res)
		return err
	}
}

func (c *MongoClient) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), 15 * time.Second)
	//return c.mongoConn.Ping(context.TODO(), nil)
	return c.mongoConn.Ping(ctx, nil)
}

func (c *MongoClient) Close() {
	ctx, _ := context.WithTimeout(context.Background(), 15 * time.Second)
	//_ = c.mongoConn.Disconnect(context.TODO())
	_ = c.mongoConn.Disconnect(ctx)
}