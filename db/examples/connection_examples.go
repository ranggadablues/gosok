package examples

import (
	"fmt"
	"log"

	"github.com/ranggadablues/gosok/db"
	"github.com/ranggadablues/gosok/db/ref"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ConnectionExamples demonstrates different connection patterns
func ConnectionExamples() {
	fmt.Println("=== MongoDB Connection Patterns ===")

	// Pattern 1: Individual instances (each service creates its own pool)
	individualInstanceExample()

	// Pattern 2: Singleton pattern (all services share one pool)
	singletonInstanceExample()
}

// individualInstanceExample shows each service creating its own connection
func individualInstanceExample() {
	fmt.Println("\n--- Pattern 1: Individual Instances ---")

	// Service 1
	service1Mongo := db.NewMongo()
	if service1Mongo == nil {
		log.Println("Service 1: Failed to connect")
		return
	}
	defer service1Mongo.Close()

	// Service 2
	service2Mongo := db.NewMongo()
	if service2Mongo == nil {
		log.Println("Service 2: Failed to connect")
		return
	}
	defer service2Mongo.Close()

	// Each service has its own connection pool
	// Total connections: 2 pools × (5-20 connections each) = 10-40 connections

	// Service 1 operations
	var users1 []bson.M
	err := service1Mongo.Find(&users1, bson.M{}, "users", ref.WithLimit(5))
	if err != nil {
		log.Printf("Service 1 find error: %v", err)
	} else {
		fmt.Printf("Service 1: Found %d users\n", len(users1))
	}

	// Service 2 operations
	var users2 []bson.M
	err = service2Mongo.Find(&users2, bson.M{}, "users", ref.WithLimit(3))
	if err != nil {
		log.Printf("Service 2 find error: %v", err)
	} else {
		fmt.Printf("Service 2: Found %d users\n", len(users2))
	}
}

// singletonInstanceExample shows services sharing one connection pool
func singletonInstanceExample() {
	fmt.Println("\n--- Pattern 2: Singleton Pattern (Recommended) ---")

	// Service 1 using singleton
	service1Mongo := GetMongoInstance()
	if service1Mongo == nil {
		log.Println("Service 1: Failed to get singleton instance")
		return
	}

	// Service 2 using same singleton instance
	service2Mongo := GetMongoInstance()
	if service2Mongo == nil {
		log.Println("Service 2: Failed to get singleton instance")
		return
	}

	// Both services share the same connection pool
	// Total connections: 1 pool × (5-20 connections) = 5-20 connections

	fmt.Printf("Service 1 and Service 2 share same instance: %t\n", service1Mongo == service2Mongo)

	// Service 1 operations
	var users1 []bson.M
	err := service1Mongo.Find(&users1, bson.M{}, "users", ref.WithLimit(5))
	if err != nil {
		log.Printf("Service 1 find error: %v", err)
	} else {
		fmt.Printf("Service 1: Found %d users\n", len(users1))
	}

	// Service 2 operations (using same connection pool)
	var users2 []bson.M
	err = service2Mongo.Find(&users2, bson.M{}, "users", ref.WithLimit(3))
	if err != nil {
		log.Printf("Service 2 find error: %v", err)
	} else {
		fmt.Printf("Service 2: Found %d users\n", len(users2))
	}

	// Close singleton when application shuts down
	// defer CloseMongoInstance() // Call this in your main() function
}

// ConnectionHealthExample demonstrates connection health checking
func ConnectionHealthExample() {
	fmt.Println("\n=== Connection Health Checking ===")

	mongo := db.NewMongo()
	if mongo == nil {
		log.Println("Failed to initialize MongoDB")
		return
	}
	defer mongo.Close()

	// The ensureConnection() method is called automatically before each operation
	// It checks if connection is alive and reconnects if needed

	fmt.Println("Performing operations with automatic connection health checks...")

	// Operation 1: Insert (will check connection health first)
	result, err := mongo.InsertOne("test", bson.M{"name": "test", "timestamp": "now"})
	if err != nil {
		log.Printf("Insert error: %v", err)
	} else {
		fmt.Printf("Inserted document with ID: %v\n", result.InsertedID)
	}

	// Operation 2: Find (will check connection health first)
	var docs []bson.M
	err = mongo.Find(&docs, bson.M{}, "test", ref.WithLimit(1))
	if err != nil {
		log.Printf("Find error: %v", err)
	} else {
		fmt.Printf("Found %d documents\n", len(docs))
	}

	// If connection was lost between operations, it would automatically reconnect
	fmt.Println("All operations completed with automatic connection management!")
}
