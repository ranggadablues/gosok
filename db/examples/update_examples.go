package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/ranggadablues/gosok/db"
	"github.com/ranggadablues/gosok/db/ref"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// UpdateExamples demonstrates various update operations including upsert functionality
func UpdateExamples() {
	// Initialize MongoDB manager
	mongoManager := db.NewMongo()
	defer mongoManager.Close()

	// Example 1: Basic UpdateOneSet without upsert (existing behavior)
	basicUpdateExample(mongoManager)

	// Example 2: UpdateOneSet with upsert
	upsertUpdateExample(mongoManager)

	// Example 3: UpdateOneSetPipeline with upsert
	upsertPipelineExample(mongoManager)

	// Example 4: UpdateManySet with upsert
	upsertManyExample(mongoManager)

	// Example 5: Complex upsert scenarios
	complexUpsertExample(mongoManager)
}

// basicUpdateExample shows the traditional update behavior (no upsert)
func basicUpdateExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Basic Update Example (No Upsert) ===")

	// First, insert a test user
	user := bson.M{
		"name":    "John Doe",
		"email":   "john.doe@example.com",
		"age":     30,
		"created": time.Now(),
		"updated": time.Now(),
	}

	insertedID, err := mongoManager.InsertOne("users", user)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return
	}
	fmt.Printf("Inserted user with ID: %v\n", insertedID)

	// Update the user (this will work because the user exists)
	updateData := bson.M{
		"age":     31,
		"status":  "active",
		"updated": time.Now(),
	}

	err = mongoManager.UpdateOneSet("users", bson.M{"email": "john.doe@example.com"}, updateData)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
	} else {
		fmt.Println("Successfully updated existing user")
	}

	// Try to update a non-existent user (this will fail without upsert)
	nonExistentUpdate := bson.M{
		"age":     25,
		"status":  "inactive",
		"updated": time.Now(),
	}

	err = mongoManager.UpdateOneSet("users", bson.M{"email": "nonexistent@example.com"}, nonExistentUpdate)
	if err != nil {
		fmt.Printf("Expected failure - user doesn't exist: %v\n", err)
	} else {
		fmt.Println("Unexpected success - this shouldn't happen without upsert")
	}
}

// upsertUpdateExample demonstrates upsert functionality with UpdateOneSet
func upsertUpdateExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Upsert Update Example ===")

	// Update with upsert - this will create the document if it doesn't exist
	upsertData := bson.M{
		"name":    "Jane Smith",
		"age":     28,
		"status":  "active",
		"created": time.Now(),
		"updated": time.Now(),
	}

	// This will create a new document because the email doesn't exist
	err := mongoManager.UpdateOneSet(
		"users",
		bson.M{"email": "jane.smith@example.com"},
		upsertData,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert user: %v", err)
	} else {
		fmt.Println("Successfully upserted user (created new document)")
	}

	// Now update the same user with different data
	updateData := bson.M{
		"age":     29,
		"status":  "premium",
		"updated": time.Now(),
	}

	err = mongoManager.UpdateOneSet(
		"users",
		bson.M{"email": "jane.smith@example.com"},
		updateData,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
	} else {
		fmt.Println("Successfully updated existing user")
	}

	// Verify the user exists and has the correct data
	var foundUser bson.M
	err = mongoManager.FindOne(&foundUser, bson.M{"email": "jane.smith@example.com"}, "users")
	if err != nil {
		log.Printf("Failed to find upserted user: %v", err)
	} else {
		fmt.Printf("Found upserted user: %s, age: %v, status: %s\n",
			foundUser["name"], foundUser["age"], foundUser["status"])
	}
}

// upsertPipelineExample demonstrates upsert with pipeline updates
func upsertPipelineExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Upsert Pipeline Example ===")

	// Create a user with pipeline update (using $set in pipeline format)
	pipelineUpdate := bson.M{
		"name":    "Bob Wilson",
		"age":     35,
		"status":  "active",
		"created": time.Now(),
		"updated": time.Now(),
	}

	err := mongoManager.UpdateOneSetPipeline(
		"users",
		bson.M{"email": "bob.wilson@example.com"},
		pipelineUpdate,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert user with pipeline: %v", err)
	} else {
		fmt.Println("Successfully upserted user with pipeline update")
	}

	// Update the same user with a more complex pipeline operation
	complexPipelineUpdate := bson.M{
		"$inc": bson.M{"age": 1},
		"$set": bson.M{
			"status":  "premium",
			"updated": time.Now(),
		},
	}

	err = mongoManager.UpdateOneSetPipeline(
		"users",
		bson.M{"email": "bob.wilson@example.com"},
		complexPipelineUpdate,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to update user with complex pipeline: %v", err)
	} else {
		fmt.Println("Successfully updated user with complex pipeline")
	}
}

// upsertManyExample demonstrates upsert with UpdateManySet
func upsertManyExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Upsert Many Example ===")

	// Update multiple users with upsert
	// Note: UpdateMany with upsert will only create one document per unique filter
	// In this case, we're updating users with different emails, so it will create multiple documents

	// First user
	err := mongoManager.UpdateManySet(
		"users",
		bson.M{"email": "alice@example.com"},
		bson.M{
			"name":    "Alice Johnson",
			"age":     25,
			"status":  "active",
			"created": time.Now(),
			"updated": time.Now(),
		},
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert Alice: %v", err)
	} else {
		fmt.Println("Successfully upserted Alice")
	}

	// Second user
	err = mongoManager.UpdateManySet(
		"users",
		bson.M{"email": "charlie@example.com"},
		bson.M{
			"name":    "Charlie Brown",
			"age":     40,
			"status":  "active",
			"created": time.Now(),
			"updated": time.Now(),
		},
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert Charlie: %v", err)
	} else {
		fmt.Println("Successfully upserted Charlie")
	}

	// Verify both users were created
	var users []bson.M
	err = mongoManager.Find(&users, bson.M{
		"email": bson.M{"$in": []string{"alice@example.com", "charlie@example.com"}},
	}, "users")
	if err != nil {
		log.Printf("Failed to find upserted users: %v", err)
	} else {
		fmt.Printf("Found %d upserted users\n", len(users))
		for _, user := range users {
			fmt.Printf("- %s (%s), age: %v\n", user["name"], user["email"], user["age"])
		}
	}
}

// complexUpsertExample demonstrates more complex upsert scenarios
func complexUpsertExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Complex Upsert Example ===")

	// Example 1: Upsert with specific ID
	userWithID := bson.M{
		"name":    "David Miller",
		"age":     32,
		"status":  "active",
		"created": time.Now(),
		"updated": time.Now(),
	}

	err := mongoManager.UpdateOneSet(
		"users",
		bson.M{"email": "david.miller@example.com"},
		userWithID,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert user: %v", err)
	} else {
		fmt.Println("Successfully upserted user with specific email")
	}

	// Example 2: Upsert with complex filter
	complexFilter := bson.M{
		"$and": []bson.M{
			{"email": "emma@example.com"},
			{"status": bson.M{"$exists": false}},
		},
	}

	complexUpdate := bson.M{
		"name":    "Emma Davis",
		"age":     27,
		"status":  "new",
		"created": time.Now(),
		"updated": time.Now(),
	}

	err = mongoManager.UpdateOneSet(
		"users",
		complexFilter,
		complexUpdate,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to upsert user with complex filter: %v", err)
	} else {
		fmt.Println("Successfully upserted user with complex filter")
	}

	// Example 3: Conditional upsert based on existing data
	conditionalUpdate := bson.M{
		"$set": bson.M{
			"name":    "Frank Wilson",
			"age":     45,
			"status":  "vip",
			"updated": time.Now(),
		},
		"$setOnInsert": bson.M{
			"created": time.Now(),
			"source":  "upsert",
		},
	}

	err = mongoManager.UpdateOneSetPipeline(
		"users",
		bson.M{"email": "frank@example.com"},
		conditionalUpdate,
		ref.WithUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to conditional upsert: %v", err)
	} else {
		fmt.Println("Successfully performed conditional upsert")
	}

	// Show all users created in this example
	var allUsers []bson.M
	err = mongoManager.Find(&allUsers, bson.M{}, "users", ref.WithLimit(10))
	if err != nil {
		log.Printf("Failed to find users: %v", err)
	} else {
		fmt.Printf("\nTotal users in database: %d\n", len(allUsers))
		for _, user := range allUsers {
			fmt.Printf("- %s (%s), age: %v, status: %s\n",
				user["name"], user["email"], user["age"], user["status"])
		}
	}
}

// CleanupExample demonstrates cleanup operations
func CleanupExample() {
	fmt.Println("\n=== Cleanup Example ===")

	mongoManager := db.NewMongo()
	defer mongoManager.Close()

	// Clean up test data
	err := mongoManager.DeleteMany("users", bson.M{})
	if err != nil {
		log.Printf("Failed to cleanup users: %v", err)
	} else {
		fmt.Println("Successfully cleaned up test data")
	}
}
