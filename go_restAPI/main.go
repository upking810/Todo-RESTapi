//A RESTful API in Go and Gin, allowing to construct API and store to databases
//This project refers the tutorial from https://thedevsaddam.medium.com/build-restful-api-service-in-golang-using-gin-gonic-framework-85b1a6e176f3 by Saddam H

package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Define the API routers
func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", CreateTodo)
		v1.GET("/", FetchAllTodo)
		v1.GET("/:id", FetchSingleTodo)
		v1.PUT("/:id", UpdateTodo)
		v1.DELETE("/:id", DeleteTodo)
	}
	router.Run()
}

var db *gorm.DB

func init() {
	//open a db connection
	var err error
	//alter this line of password and database name to your own
	//db, err = gorm.Open("mysql", "databaseName:password@/schema_name?charset=utf8&parseTime=True&loc=Local")
	db, err = gorm.Open("mysql", "root:password@/go_restful_api?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	//Gorm can migrate facilities. Migrate the schema
	db.AutoMigrate(&TodoModel{})
}

// construct the todoModel struct and the transformedTodo struct
type (
	// TodoModel describes a todoModel type
	TodoModel struct {
		gorm.Model        //embedds ID, CreatedAt, UpdatedAt, DeletedAt
		Title      string `json:"title"`
		Completed  int    `json:"completed"`
	}

	// transformedTodo represents a formatted todo. This does not includes the CreatedAt, UpdatedAt, DeletedAt
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

// CreateTodo add a new todo
func CreateTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	//add the todo model
	todo := TodoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	//send the response
	c.JSON(201, gin.H{"status": "success", "message": "Todo item created successfully!", "resourceId": todo.ID})
}

// FetchAllTodo fetch all todos
func FetchAllTodo(c *gin.Context) {
	var todos []TodoModel
	var _todos []transformedTodo
	//fetch all todos
	db.Find(&todos)
	if len(todos) <= 0 {
		c.JSON(404, gin.H{"status": "not found", "message": "No todo found!"})
		return
	}
	//transforms the todos for building a good response
	//read the todos array, check the completion, assign it to the transformedTodo array
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(200, gin.H{"status": "success", "data": _todos})
}

// FetchSingleTodo fetch a single todo
func FetchSingleTodo(c *gin.Context) {
	var todo TodoModel
	todoID := c.Param("id") //get the id from the url
	db.First(&todo, todoID) //fetch the database record by ID
	//if the todo is not found
	if todo.ID == 0 {
		c.JSON(404, gin.H{"status": "not found", "message": "No todo found!"})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	// find the todo and return the data
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(200, gin.H{"status": "success", "data": _todo})
}

// UpdateTodo update a todo
func UpdateTodo(c *gin.Context) {
	var todo TodoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(404, gin.H{"status": "not found", "message": "No todo found!"})
		return
	}
	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(200, gin.H{"status": "success", "message": "Todo updated successfully!"})
}

// DeleteTodo remove a todo
func DeleteTodo(c *gin.Context) {
	//get the id from the url
	todoID := c.Param("id")
	var todo TodoModel
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(404, gin.H{"status": "not found", "message": "No todo found!"})
		return
	}
	db.Delete(&todo)
	c.JSON(200, gin.H{"status": "success", "message": "Todo deleted successfully!"})
}
