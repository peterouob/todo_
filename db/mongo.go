package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var Mgo *mongo.Client

func InitMongo() {
	opts := options.Client().ApplyURI(os.Getenv("MONGO"))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Panicf("error in connect mongo db:%s", err.Error())
		return
	}
	Mgo = client
}
