package driver

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client *mongo.Client
}

var Mongo MongoDB

func ConnectMongoDb(userName, password string) {
	if Mongo.Client != nil {
		return
	}
	connStr := fmt.Sprintf("mongodb://%v:%v@localhost:27017/", userName, password)
	// Connect to MongoDB
	client, err := mongo.NewClient((options.Client().ApplyURI(connStr)))
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	fmt.Println("Connect Successfully!")
	Mongo.Client = client
}

func GetMongoDb() (*MongoDB, error) {
	mongo := MongoDB{}
	if Mongo.Client == nil {
		return &mongo, nil
	}
	return &Mongo, nil
}
