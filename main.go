package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kingcortex/cinema-api/config"
	"github.com/kingcortex/cinema-api/middleware"
	"github.com/kingcortex/cinema-api/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func initialize() {
	log.Println("Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	log.Println("Environment variables loaded.")

	log.Println("Setting application configuration...")
	config.MyConfig = &config.Config{
		MONGO_URI:             os.Getenv("MONGO_URI"),
		MOVIE_COLLECTION_NAME: os.Getenv("MOVIE_COLLECTION_NAME"),
		USER_COLLECTION_NAME:  os.Getenv("USER_COLLECTION_NAME"),
		DB_NAME:               os.Getenv("DB_NAME"),
		JWT_SECRET:            os.Getenv("JWT_SECRET"),
	}
	log.Println("Application configuration set:", config.MyConfig)

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)

	log.Println("Connecting to MongoDB...")
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB. Pinging database...")
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	log.Println("Ping successful. Connected to MongoDB deployment!")
	mongoClient = client
}

func main() {
	defer func() {
		log.Println("Disconnecting from MongoDB...")
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		log.Println("Disconnected from MongoDB.")
	}()
	initialize()
	movieCollection := mongoClient.Database(config.MyConfig.DB_NAME).Collection(config.MyConfig.MOVIE_COLLECTION_NAME)
	userCollerction := mongoClient.Database(config.MyConfig.DB_NAME).Collection(config.MyConfig.USER_COLLECTION_NAME)

	movieService := services.MovieService{MongoCollection: *movieCollection}
	authService := services.AuthService{MongoCollection: userCollerction}

	r := gin.Default()

	protected := r.Group("/api", middleware.AuthMiddleware(userCollerction))

	protected.GET("/movies", movieService.GetAll)
	protected.POST("/movies", movieService.CreateOneMovie)
	protected.PUT("/movies/:id", movieService.UpdateMovieByID)
	protected.GET("/movies/:id", movieService.GetOneMovie)
	r.POST("/register", authService.Register)
	r.POST("/login", authService.Login)

	r.Run(":8080")
}
