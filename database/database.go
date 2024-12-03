package database

import (
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var db *mongo.Database

func Init() bool {
	uri := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetTLSConfig(&tls.Config{})
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println(err)
		return false
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Unable to connect to DB: " + err.Error())
		return false
	}

	log.Printf("Connected to MongoDB URI", uri)

	dbs, err := client.ListDatabaseNames(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		log.Fatal("Unable to connect to DB")
		return false
	}

	for i := range dbs {
		if dbs[i] == os.Getenv("MONGODB_DATABASE") {
			log.Println("Found " + os.Getenv("MONGODB_DATABASE") + " DB")
		}
	}

	db = client.Database(os.Getenv("MONGODB_DATABASE"))

	return true
}

func Find(collection string, filters bson.M, opt *options.FindOptions) *mongo.Cursor {
	option := options.Find()
	if opt != nil {
		option = opt
	}

	cur, err := db.Collection(collection).Find(context.Background(), filters, option)

	if err != nil {
		log.Println(err)
		return nil
	}

	return cur
}

func FindOne(collection string, filters bson.M, opt *options.FindOneOptions) *mongo.SingleResult {
	if opt == nil {
		opt = options.FindOne()
	}

	cur := db.Collection(collection).FindOne(context.Background(), filters, opt)

	if cur == nil {
		return nil
	}

	return cur
}

func Count(collection string, filters bson.M) int64 {
	count, err := db.Collection(collection).CountDocuments(context.Background(), filters)

	if err != nil {
		log.Println(err)
		return -1
	}

	return count
}

func InsertOne(collection string, object interface{}) (*mongo.InsertOneResult, error) {
	res, err := db.Collection(collection).InsertOne(context.Background(), object, options.InsertOne())

	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateOne(collection string, filters bson.M, object interface{}) (*mongo.UpdateResult, error) {
	data := bson.M{"$set": object}

	res, err := db.Collection(collection).UpdateOne(context.Background(), filters, data, options.Update())

	if err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteOne(collection string, filter bson.M) (*mongo.DeleteResult, error) {
	res, err := db.Collection(collection).DeleteOne(context.Background(), filter, options.Delete())

	if err != nil {
		return nil, err
	}

	return res, nil
}
