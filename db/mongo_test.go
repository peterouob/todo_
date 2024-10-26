package db

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

func init() {
	if err := godotenv.Load("test.env"); err != nil {
		log.Panicf("error to load env file ... :%s", err.Error())
	}
}

func TestInitMongo(t *testing.T) {
	opts := options.Client().ApplyURI(os.Getenv("MONGO"))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		t.Logf("error in connect mogno:%s", err.Error())
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			t.Errorf("error in close mongo:%s", err.Error())
		}
	}()
}
