package lr

import (
	"context"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbName = "LoginRadiusDev"
)

var mongoClientOnce sync.Once
var mongoClient *mongo.Database

func connectToMongo(connStr string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(uint64(1000000000)))
	defer cancel()

	clientOptions := options.Client().ApplyURI(connStr)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return client
}

func NewMongoClient() *mongo.Database {
	devMongoDb := os.Getenv("MONGO_URI")
	mongoClientOnce.Do(func() {
		mongoClient = connectToMongo(devMongoDb).Database(dbName)
	})

	return mongoClient
}

func CreateAppidToOrgidMapping(mongoClient *mongo.Database, appid int, orgid string) error {

	if _, err := mongoClient.Collection("Org_App").InsertOne(context.Background(), bson.M{"appid": appid, "orgid": orgid}); err != nil {
		return err
	}

	return nil
}

func GetAppIdFromOrgIdMapping(mongoClient *mongo.Database, orgid string) (int, error) {
	var result struct {
		Appid int `bson:"appid"`
	}

	if err := mongoClient.Collection("Org_App").FindOne(context.Background(), bson.M{"orgid": orgid}).Decode(&result); err != nil {
		return 0, err
	}

	return result.Appid, nil
}
