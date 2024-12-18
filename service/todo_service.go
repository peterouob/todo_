package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/peterouob/todo_/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strconv"
)

func GetAllTodo(c *gin.Context) {
	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	data, err := findAllTodo(context.TODO(), userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  data,
		})
		return
	}
	if len(data) == 0 {
		data = []model.Todo{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  data,
	})
}

func GetTodoFilterDone(c *gin.Context) {
	var req model.FindTodoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	data, err := findTodoFilterDone(req.Done, userId, context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  data,
	})
}

func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	data, err := findById(objectId, userId, context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  data,
	})
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	if err := deleteTodo(objectId, userId, context.TODO()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"delete id": id,
	})
}

func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	req := model.UpdateTodoRequest{}

	if err := c.ShouldBind(&req); !errors.Is(err, io.EOF) && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	if err := updateTodo(objectId, req, userId, context.TODO()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func CreateTodo(c *gin.Context) {
	todo := model.NewTodo()
	var err error
	if err = c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	todo_id, err := createTodo(todo, userId, context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	todo.Id = todo_id.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  todo,
	})
}

func Done(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"id":   id,
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	if err := doneTodo(objectId, userId, context.TODO()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"id":   id,
	})
}

func FilterByMonthAndYear(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("m"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "error in format query string",
		})
		return
	}
	year, err := strconv.Atoi(c.Query("y"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "error in format query string",
		})
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "invalid user ID in cookie",
		})
		return
	}

	data, err := findByMonthAndYear(month, year, userId, context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}
