package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type task struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var tasks = []task{
	{Id: 1, Title: "title 1", Description: "Description 1", Price: 1},
	{Id: 2, Title: "title 2", Description: "Description 2", Price: 2.8},
	{Id: 3, Title: "title 3", Description: "Description 3", Price: 178.5},
}

func main() {



	r := gin.Default()
	r.GET("/ping", ping)
	r.GET("/", start)
	r.GET("/tasks",getTasks)

	r.POST("/task",addTask)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func start(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
} 

func ping(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, tasks)
}

func addTask(context *gin.Context) {
	var newTask task
	if err:= context.BindJSON(&newTask) ; err !=nil{
		return
	}
	fmt.Println(newTask)
	tasks = append(tasks, newTask)
	context.IndentedJSON(http.StatusCreated, newTask)
}
