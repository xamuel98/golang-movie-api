package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}

	// Using context.Background() as the root of the context tree
	/**
	* In Go, the context.Background() function returns a non-nil, empty Context.
	* It is often used as the root of a context tree and serves as a starting point
	* for creating more specific context objects using the context.WithXXX functions.
	*
	* The context package in Go provides a way to carry deadlines, cancellation signals,
	* and other request-scoped values across API boundaries and between processes.
	* It is commonly used in concurrent and distributed systems to manage the lifecycle of
	* operations, propagate deadlines, and handle cancellations.
	*
	* The context.Background() function is used when you don't have a more specific Context
	* available or when you are creating the root of a context tree. It creates a basic,
	* empty Context without any associated values, deadlines, or cancellation signals.
	 */
	rootContext := context.Background()

	client, err := mongo.Connect(rootContext, options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	err := client.Disconnect(rootContext)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()

	if err := client.Ping(rootContext, nil); err != nil {
		log.Fatal("Error pinging MongoDB server: ", err)
		panic(err)
	}

	fmt.Printf("Pinged your deployment. You successfully connected to MongoDB!\n")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster1").Collection(collectionName)
	return collection
}
