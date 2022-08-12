package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type task struct {
	ID          string  `json:"_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var collection *mongo.Collection

func main() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("go_rest").Collection("trainers")

	r := gin.Default()
	r.GET("/ping", ping)

	r.GET("/task/:id", getTask)
	r.GET("/tasks/", getTasks)
	r.POST("/task/", addTask)
	r.DELETE("/task/:id", deleteTask)
	r.PUT("/task/:id", updateTask)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getTasks(c *gin.Context) {

	//filter := bson.D{{"_id", bson.M{"$eq": bsontype.Null}}}

	cursor, err := collection.Find(context.TODO(), bson.D{})
	var results []bson.M
	if err != nil {
		log.Fatal(err)
	}else {
		for cursor.Next(c) {
			// Declare a result BSON object
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}else {
				results = append(results,result)
			}
		}
	}
	c.IndentedJSON(http.StatusCreated, results)

}

func deleteTask(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idParam)
	filter := bson.D{{"_id", id}}

	result, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, result)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func getTask(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idParam)
	filter := bson.D{{"_id", id}}

	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return
		}
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusCreated, result)
}

func addTask(c *gin.Context) {
	var newTask task
	if err := c.BindJSON(&newTask); err != nil {
		return
	}
	insertResult, err := collection.InsertOne(context.TODO(), newTask)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, insertResult)
}

func updateTask(c *gin.Context) {
	var newTask task

	if err := c.BindJSON(&newTask); err != nil {
		return
	}

	idParam := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	filter := bson.D{{"_id", id}}

	update := bson.D{
		{"$set", bson.D{{"title", newTask.Title}}},
		{"$set", bson.D{{"description", newTask.Description}}},
		{"$set", bson.D{{"price", newTask.Price}}},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusCreated, result)
}
