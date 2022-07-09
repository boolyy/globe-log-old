package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"context"

	"github.com/boolyy/globe-log/go-backend/controllers"
	"github.com/boolyy/globe-log/go-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Listening on port:3000")
// }

func loadMongoUri() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	mongoURI := os.Getenv("ATLAS_URI")
	if mongoURI == "" {
		return "", fmt.Errorf("uri environment variable not found")
	}

	return os.Getenv("ATLAS_URI"), nil

}

// func ConnectMongo() *mongo.Client {
// 	mongoURI, err := loadMongoUri()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

// 	err = client.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		log.Fatal("error while trying to ping mongo", err)
// 	}

// 	defer func() {
// 		if err = client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	return client
// }

func main() {
	// http.Get("/hi")
	// http.HandleFunc("/", hello)
	// http.ListenAndServe(":3000", nil)

	mongoURI, err := loadMongoUri()
	if err != nil {
		log.Fatal(err)
	}

	mongoContext := context.Background()

	client, err := mongo.Connect(mongoContext, options.Client().ApplyURI(mongoURI))

	err = client.Ping(mongoContext, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	defer func() {
		if err = client.Disconnect(mongoContext); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//authController := controllers.AuthController{UserCollection: client.Database("globe_data").Collection("users"), MongoContext: mongoContext}
	userCollection := client.Database("globe_data").Collection("users")
	userService := services.NewUserService(userCollection, mongoContext)

	controller := controllers.Controller{UserService: userService}

	basePath := router.Group("/")
	controller.RegisterAuthRoutes(basePath)
	controller.RegisterLocationRoutes(basePath)
	controller.RegisterProfileRoutes(basePath)

	router.Run(":3000")

}
