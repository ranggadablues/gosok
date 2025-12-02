package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/ranggadablues/gosok/db"
	"github.com/ranggadablues/gosok/db/ref"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	testUserEmail = "john@example.com"
)

// Example demonstrates how to use the MongoDB helper
func Example() {
	basicExample()
	advancedExample()
	UpdateExamples()
	CleanupExample()
}

// basicExample shows basic MongoDB operations
func basicExample() {
	// Initialize MongoDB manager
	mongoManager := db.NewMongo()
	defer mongoManager.Close()

	// Example 1: Get a client directly
	defaultClient := mongoManager.GetClient()
	fmt.Println("Got default client:", defaultClient != nil)

	// Example 2: Get a collection
	usersCollection := mongoManager.GetCollection("users")
	fmt.Println("Got users collection:", usersCollection != nil)

}

// advancedExample shows more complex MongoDB operations
func advancedExample() {
	// Initialize MongoDB manager
	mongoManager := db.NewMongo()
	defer mongoManager.Close()

	// Example 3: Insert a document using the simplified API
	user := bson.M{
		"name":    "John Doe",
		"email":   testUserEmail,
		"created": time.Now(),
	}

	insertedID, err := mongoManager.InsertOne("users", user)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
	} else {
		fmt.Printf("Inserted user with ID: %v\n", insertedID)
	}

	// Example 4: Insert multiple documents
	users := []any{
		bson.M{"name": "Alice", "email": "alice@example.com", "age": 25},
		bson.M{"name": "Bob", "email": "bob@example.com", "age": 30},
		bson.M{"name": "Charlie", "email": "charlie@example.com", "age": 35},
	}

	insertedIDs, err := mongoManager.InsertMany("users", users)
	if err != nil {
		log.Printf("Failed to insert users: %v", err)
	} else {
		fmt.Printf("Inserted %d users\n", len(insertedIDs))
	}

	// Example 5: Find one document
	var foundUser bson.M
	err = mongoManager.FindOne(&foundUser, bson.M{"email": testUserEmail}, "users")
	if err != nil {
		log.Printf("Failed to find user: %v", err)
	} else {
		fmt.Printf("Found user: %v\n", foundUser["name"])
	}

	// Example 6: Find with sort and limit
	var sortedUsers []bson.M
	err = mongoManager.Find(&sortedUsers, bson.M{}, "users",
		ref.WithLimit(5),
		ref.WithSort(bson.D{{Key: "created", Value: -1}}),
	)
	if err != nil {
		log.Printf("Failed to find sorted users: %v", err)
	} else {
		fmt.Printf("Found %d sorted users\n", len(sortedUsers))
	}

	// Example 7: Find with projection
	var userEmails []bson.M
	err = mongoManager.Find(&userEmails, bson.M{}, "users",
		ref.WithProjection(bson.D{{Key: "email", Value: 1}, {Key: "_id", Value: 0}}),
		ref.WithLimit(3),
	)
	if err != nil {
		log.Printf("Failed to find user emails: %v", err)
	} else {
		fmt.Printf("Found %d user emails:\n", len(userEmails))
		for _, user := range userEmails {
			fmt.Printf("- %s\n", user["email"])
		}
	}

	// Example 8: Delete operations
	err = mongoManager.DeleteOne("users", bson.M{"email": testUserEmail})
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
	}

	// Example 9: Advanced Find examples
	showAdvancedFindExamples(mongoManager)
}

// showAdvancedFindExamples demonstrates advanced find operations
func showAdvancedFindExamples(mongoManager db.IMongoLib) {
	// Example 1: Basic Find without options
	var users []bson.M
	err := mongoManager.Find(&users, bson.M{}, "users")
	if err != nil {
		log.Printf("Failed to find users: %v", err)
		return
	}
	fmt.Printf("Found %d users\n", len(users))

	// Example 2: Find with limit and sort
	var recentUsers []bson.M
	err = mongoManager.Find(&recentUsers, bson.M{}, "users",
		ref.WithLimit(10),
		ref.WithSort(bson.D{{Key: "created", Value: -1}}),
	)
	if err != nil {
		log.Printf("Failed to find recent users: %v", err)
	} else {
		fmt.Printf("Found %d recent users\n", len(recentUsers))
	}

	// Example 3: Find with pagination (skip and limit)
	var pagedUsers []bson.M
	err = mongoManager.Find(&pagedUsers, bson.M{}, "users",
		ref.WithSkip(20),
		ref.WithLimit(10),
		ref.WithSort(bson.D{{Key: "name", Value: 1}}),
	)
	if err != nil {
		log.Printf("Failed to find paged users: %v", err)
	} else {
		fmt.Printf("Found %d users on page 3\n", len(pagedUsers))
	}

	// Example 4: Find with projection (only specific fields)
	var userNames []bson.M
	err = mongoManager.Find(&userNames, bson.M{}, "users",
		ref.WithProjection(bson.D{{Key: "name", Value: 1}, {Key: "email", Value: 1}, {Key: "_id", Value: 0}}),
		ref.WithLimit(5),
	)
	if err != nil {
		log.Printf("Failed to find user names: %v", err)
	} else {
		fmt.Printf("Found %d user names and emails\n", len(userNames))
		for _, user := range userNames {
			fmt.Printf("- %s (%s)\n", user["name"], user["email"])
		}
	}

	// Example 5: Complex query with all options
	var complexResult []bson.M
	err = mongoManager.Find(&complexResult,
		bson.M{"age": bson.M{"$gte": 18, "$lt": 65}},
		"users",
		ref.WithSort(bson.D{{Key: "age", Value: 1}, {Key: "name", Value: 1}}),
		ref.WithSkip(0),
		ref.WithLimit(25),
		ref.WithProjection(bson.D{
			{Key: "name", Value: 1},
			{Key: "age", Value: 1},
			{Key: "email", Value: 1},
			{Key: "_id", Value: 0},
		}),
	)
	if err != nil {
		log.Printf("Failed to find adult users: %v", err)
	} else {
		fmt.Printf("Found %d adult users sorted by age\n", len(complexResult))
	}

	// Example 6: Delete many documents
	err = mongoManager.DeleteMany("users", bson.M{"age": bson.M{"$lt": 18}})
	if err != nil {
		log.Printf("Failed to delete users: %v", err)
	}
}
