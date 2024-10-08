package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos []Todo = []Todo{
	{Id: 1, Title: "Do Laundry", Done: false},
	{Id: 2, Title: "Buy Groceries", Done: true},
	{Id: 3, Title: "Clean House", Done: false},
}

func getTodos(context *gin.Context) {
	context.JSON(http.StatusOK, todos)
}

func searchTodo(id uint) (*Todo, error) {
	for i, todo := range todos {
		if id == todo.Id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("no matching todo")
}

func getTodoByID(context *gin.Context) {
	id := context.Param("id")

	if n, err := strconv.Atoi(id); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "id could not be parsed to uint! id: " + string(id)})
	} else {
		if todo, err := searchTodo(uint(n)); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "could not find todo with given id: " + string(n)})
		} else {
			context.JSON(http.StatusOK, todo)
		}
	}

}

func addTodo(context *gin.Context) {
	var newTodo Todo
	if err := context.BindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		todos = append(todos, newTodo)
		context.JSON(http.StatusOK, todos)
	}
}

func updateTodo(context *gin.Context) {
	var updateTodo Todo
	if err := context.BindJSON(&updateTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if todo, err := searchTodo(updateTodo.Id); err == nil {
		*todo = updateTodo
		context.JSON(http.StatusOK, *todo)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteTodo(context *gin.Context) {
	id := context.Param("id")

	if n, err := strconv.Atoi(id); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "id could not be parsed to uint! id: " + string(id)})
	} else {
		if todo, err := searchTodo(uint(n)); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "could not find todo with given id: " + string(n)})
		} else {
			*todo = todos[len(todos)-1]
			todos = todos[:len(todos)-1]
			context.JSON(http.StatusOK, todos)
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodoByID)
	router.POST("/todos", addTodo)
	router.PUT("/todos", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run(":8080")
}
