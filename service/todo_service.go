package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/peterouob/todo_/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

func GetAllTodo(c *gin.Context) {
	data, err := findAllTodo(context.TODO())
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

func GetTodoFilterDone(c *gin.Context) {
	var req model.FindTodoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	data, err := findTodoFilterDone(req.Done, context.TODO())
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
	data, err := findById(objectId, context.TODO())
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
	if err := deleteTodo(objectId, context.TODO()); err != nil {
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
	req := model.UpdateTodoRequest{}

	if err := c.ShouldBind(&req); !errors.Is(err, io.EOF) && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	if err := updateTodo(objectId, req, context.TODO()); err != nil {
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
	if err = createTodo(todo, context.TODO()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

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
	if err := doneTodo(objectId, context.TODO()); err != nil {
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
