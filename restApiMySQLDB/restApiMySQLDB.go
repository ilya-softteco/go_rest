package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

type Task struct {
	Id       string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func dsn(viper *viper.Viper) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		viper.Get("user"), viper.Get("password"),
		viper.Get("hostDB"), viper.Get("portDB"), viper.Get("dbName"))
}

var db *sql.DB

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("mysql", dsn(viper.GetViper()))

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	r := gin.Default()
	r.GET("/ping", ping)

	r.GET("/task/:id", getTask)
	r.GET("/tasks/", getTasks)
	r.POST("/task/", addTask)
	r.DELETE("/task/:id", deleteTask)
	r.PUT("/task/:id", updateTask)

	r.Run(fmt.Sprintf(
		"%s:%s",
		viper.Get("hostServer"), viper.Get("portServer")))
}

func getTasks(c *gin.Context) {

	res, err :=db.Query("SELECT * FROM tasks")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	var tasks []Task
	for res.Next() {
		var task Task
		if err := res.Scan(&task.Id,&task.Title, &task.Description,&task.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())

		}
		tasks = append(tasks, task)
	}
	c.IndentedJSON(http.StatusOK, tasks)

}

func deleteTask(c *gin.Context) {
	idParam := c.Param("id")

	_, err :=db.Query("DELETE FROM  tasks WHERE id = ?", idParam)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, true)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func getTask(c *gin.Context) {

	idParam := c.Param("id")

	res, err :=db.Query("SELECT * FROM tasks WHERE id = ?", idParam)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	var task Task
	res.Next()
	if err := res.Scan(&task.Id,&task.Title, &task.Description,&task.Price); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, task)

}

func addTask(c *gin.Context) {
	var newTask Task
	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	sql := fmt.Sprintf(
		"INSERT INTO tasks (title, description, price) VALUES ('%s','%s',%f)",
		newTask.Title, newTask.Description, newTask.Price)
	res, err := db.Exec(sql)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, res)
}

func updateTask(c *gin.Context) {

	var newTask Task
	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	sql := fmt.Sprintf(
		"UPDATE  tasks SET title = '%s' , description = '%s' , price = %f ",
		newTask.Title, newTask.Description, newTask.Price)
	res, err := db.Exec(sql)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, res)
}
