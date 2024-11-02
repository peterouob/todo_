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
	"time"
)

func findAllTodo(ctx context.Context, userId int64) ([]model.Todo, error) {
	var result []model.Todo

	filter := bson.M{"userID": userId}
	sortFilter := bson.D{bson.E{Key: "deadTime", Value: 1}}
	opts := options.Find().SetSort(sortFilter)
	cursor, err := db.Mgo.Find(ctx, filter, opts)
	if err != nil {
		return result, errors.New("error in filter the data from collection:" + err.Error())
	}
	if err := cursor.All(ctx, &result); err != nil {
		return result, errors.New("didn't have any todo")
	}
	return result, nil
}

func findTodoFilterDone(done bool, userId int64, ctx context.Context) ([]model.Todo, error) {
	filter := bson.M{"done": done, "userID": userId}
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

func findById(id primitive.ObjectID, userId int64, ctx context.Context) (model.Todo, error) {
	var todo model.Todo

	if id.IsZero() {
		return todo, errors.New("the id cannot be empty")
	}

	filter := bson.M{"_id": bson.M{"$eq": id}, "userID": userId}
	err := db.Mgo.FindOne(ctx, filter).Decode(&todo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return todo, errors.New("not found the id which want to find")
	}
	if err != nil {
		return model.Todo{}, errors.New("error in find data from collection by id:" + err.Error())
	}
	return todo, nil
}

func deleteTodo(id primitive.ObjectID, userId int64, ctx context.Context) error {
	if id.IsZero() {
		return errors.New("the id cannot be empty")
	}
	filter := bson.M{"_id": bson.M{"$eq": id}, "userID": userId}
	_, err := db.Mgo.DeleteOne(ctx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("not found the id which want to delete")
	}
	if err != nil {
		return errors.New("error in delete data from collection: " + err.Error())
	}
	return nil
}

func updateTodo(id primitive.ObjectID, req model.UpdateTodoRequest, userId int64, ctx context.Context) error {
	if id.IsZero() {
		return errors.New("the id cannot be empty")
	}
	filter := bson.M{"_id": bson.M{"$eq": id}, "userID": userId}
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

func createTodo(todo model.Todo, id int64, ctx context.Context) (interface{}, error) {
	todo.UserID = id
	result, err := db.Mgo.InsertOne(ctx, todo)
	if err != nil {
		return nil, errors.New("error in insert data from collection:" + err.Error())
	}
	return result.InsertedID, nil
}

func doneTodo(id primitive.ObjectID, userId int64, ctx context.Context) error {
	if id.IsZero() {
		return errors.New("the id cannot be empty")
	}

	update := bson.M{"$set": bson.M{"done": true}}
	filter := bson.M{"_id": id, "userID": userId}
	result, err := db.Mgo.UpdateOne(ctx, filter, update)

	if result.MatchedCount != 1 {
		return errors.New("not found the id which want to update")
	}
	if err != nil {
		return errors.New("error in update todo :" + err.Error())
	}
	return nil
}

func findByMonthAndYear(month, year int, userId int64, ctx context.Context) ([]model.TodoGroup, error) {
	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"deadTime", bson.D{
					{"$gte", startTime},
					{"$lt", endTime},
				}},
				{"userID", userId},
			}},
		},
		{
			{"$group", bson.D{
				{"_id", bson.D{
					{"year", bson.D{{"$year", "$deadTime"}}},
					{"month", bson.D{{"$month", "$deadTime"}}},
				}},
				{"todos", bson.D{{"$push", bson.D{
					{"_id", "$_id"},
					{"title", "$title"},
					{"content", "$content"},
					{"deadTime", "$deadTime"},
					{"done", "$done"},
				}}}},
			}},
		},
	}
	cursor, err := db.Mgo.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("error in mongo aggregate pipeline :" + err.Error())
	}
	var result []model.TodoGroup
	if err := cursor.All(ctx, &result); err != nil {
		return nil, errors.New("error in get the value for result :" + err.Error())
	}
	return result, nil
}
