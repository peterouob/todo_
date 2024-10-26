package service

import (
	"context"
	"errors"
	"github.com/peterouob/todo_/db"
	"github.com/peterouob/todo_/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findAllTodo(ctx context.Context) ([]model.Todo, error) {
	filter := bson.M{}
	sortFilter := bson.D{bson.E{Key: "deadTime", Value: 1}}
	opts := options.Find().SetSort(sortFilter)
	cursor, err := db.Mgo.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.New("error in filter the data from collection:" + err.Error())
	}
	var result []model.Todo
	if err := cursor.All(ctx, &result); err != nil {
		return nil, errors.New("error in find all data from collection:" + err.Error())
	}

	return result, nil
}

func findTodoFilterDone(done bool, ctx context.Context) ([]model.Todo, error) {
	filter := bson.M{"done": done}
	sortFilter := bson.D{bson.E{Key: "deadTime", Value: 1}}

	opts := options.Find().SetSort(sortFilter)
	cursor, err := db.Mgo.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.New("error in filter the data from collection:" + err.Error())
	}
	var result []model.Todo
	if err := cursor.All(ctx, &result); err != nil {
		return nil, errors.New("error in find all data from collection:" + err.Error())
	}

	return result, nil
}

func findById(id primitive.ObjectID, ctx context.Context) (model.Todo, error) {
	var todo model.Todo

	if id.IsZero() {
		return todo, errors.New("the id cannot be empty")
	}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	err := db.Mgo.FindOne(ctx, filter).Decode(&todo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return todo, errors.New("not found the id which want to find")
	}
	if err != nil {
		return model.Todo{}, errors.New("error in find data from collection by id:" + err.Error())
	}
	return todo, nil
}

func deleteTodo(id primitive.ObjectID, ctx context.Context) error {
	if id.IsZero() {
		return errors.New("the id cannot be empty")
	}
	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err := db.Mgo.DeleteOne(ctx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("not found the id which want to delete")
	}
	if err != nil {
		return errors.New("error in delete data from collection: " + err.Error())
	}
	return nil
}
func updateTodo(id primitive.ObjectID, req model.UpdateTodoRequest, ctx context.Context) error {
	if id.IsZero() {
		return errors.New("the id cannot be empty")
	}
	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": req}
	result, err := db.Mgo.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("not found the id which want to update")
	}
	if err != nil {
		return errors.New("error in update todo :" + err.Error())
	}
	return nil
}

func createTodo(todo model.Todo, ctx context.Context) error {
	_, err := db.Mgo.InsertOne(ctx, todo)
	if err != nil {
		return errors.New("error in insert data from collection:" + err.Error())
	}
	return nil
}

func doneTodo(id primitive.ObjectID, ctx context.Context) error {
	if id.IsZero() {
		return errors.New("the id cannot be empty")

	}
	update := bson.M{"$set": bson.M{"done": true}}
	result, err := db.Mgo.UpdateByID(ctx, id, update)
	if result.MatchedCount != 1 {
		return errors.New("not found the id which want to update")
	}
	if err != nil {
		return errors.New("error in update todo :" + err.Error())
	}
	return nil
}
