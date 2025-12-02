package examples

import (
	"fmt"
	"log"

	"github.com/ranggadablues/gosok/db"
	"github.com/ranggadablues/gosok/db/ref"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	fmt.Println("Testing upsert functionality...")

	// Initialize MongoDB manager
	mongoManager := db.NewMongo()
	defer mongoManager.Close()

	// Test basic upsert
	fmt.Println("\n=== Testing Basic Upsert ===")

	// This should create a new document since it doesn't exist
	updateData := bson.M{
		"name":   "Test User",
		"email":  "test@example.com",
		"age":    25,
		"status": "active",
	}

	err := mongoManager.UpdateOneSet(
		"test_users",
		bson.M{"email": "test@example.com"},
		updateData,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert user: %v", err)
	} else {
		fmt.Println("Successfully upserted user (created new document)")
	}

	// Test updating existing document
	fmt.Println("\n=== Testing Update Existing Document ===")

	updateData = bson.M{
		"age":    26,
		"status": "premium",
	}

	err = mongoManager.UpdateOneSet(
		"test_users",
		bson.M{"email": "test@example.com"},
		updateData,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
	} else {
		fmt.Println("Successfully updated existing user")
	}

	// Verify the document exists
	fmt.Println("\n=== Verifying Document ===")
	var foundUser bson.M
	err = mongoManager.FindOne(&foundUser, bson.M{"email": "test@example.com"}, "test_users")
	if err != nil {
		log.Printf("Failed to find user: %v", err)
	} else {
		fmt.Printf("Found user: %s, age: %v, status: %s\n",
			foundUser["name"], foundUser["age"], foundUser["status"])
	}

	// Clean up
	fmt.Println("\n=== Cleaning Up ===")
	err = mongoManager.DeleteOne("test_users", bson.M{"email": "test@example.com"})
	if err != nil {
		log.Printf("Failed to delete test user: %v", err)
	} else {
		fmt.Println("Successfully cleaned up test data")
	}

	fmt.Println("\nUpsert functionality test completed!")
}
