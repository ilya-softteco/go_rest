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

type task struct {
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

}

func deleteTask(c *gin.Context) {

}

func ping(c *gin.Context) {

}
func getTask(c *gin.Context) {

	idParam := c.Param("id")

	/*sql := fmt.Sprintf(
		"SELECT * FROM  tasks WHERE id = %d" ,
		idParam)*/


	res, err :=db.Query("SELECT * FROM tasks WHERE id = ?", idParam)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, res.Next())

}

func addTask(c *gin.Context) {

	var newTask task
	if err := c.BindJSON(&newTask); err != nil {
		return
	}

	sql := fmt.Sprintf(
		"INSERT INTO tasks (title, description, price) VALUES ('%s','%s',%f)",
		newTask.Title, newTask.Description, newTask.Price)
	res, err := db.Exec(sql)

	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusCreated, res)
}

func updateTask(c *gin.Context) {

}
