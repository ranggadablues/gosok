package examples

import (
	"fmt"
	"log"

	"github.com/ranggadablues/gosok/db"
	"github.com/ranggadablues/gosok/db/ref"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Constants used across all examples
const (
	errMessage   = "Error"
	printMessage = "%s: %v"
)

// basicFindExample demonstrates basic Find operation
func basicFindExample(mongoManager db.IMongoLib) {
	fmt.Println("=== Example 1: Basic Find ===")
	var allUsers []bson.M
	err := mongoManager.Find(&allUsers, bson.M{}, "users")
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d users\n", len(allUsers))
	}
}

// findWithLimitExample demonstrates Find with limit option
func findWithLimitExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 2: Find with Limit ===")
	var limitedUsers []bson.M
	err := mongoManager.Find(&limitedUsers, bson.M{}, "users", ref.WithLimit(5))
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d users (limited to 5)\n", len(limitedUsers))
	}
}

// findWithSortExample demonstrates Find with sort option
func findWithSortExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 3: Find with Sort ===")
	var sortedUsers []bson.M
	err := mongoManager.Find(&sortedUsers, bson.M{}, "users",
		ref.WithSort(bson.D{{Key: "created", Value: -1}}), // Sort by created date, newest first
		ref.WithLimit(3),
	)
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d users sorted by creation date\n", len(sortedUsers))
	}
}

// findWithPaginationExample demonstrates Find with skip and limit for pagination
func findWithPaginationExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 4: Find with Skip (Pagination) ===")
	var page2Users []bson.M
	err := mongoManager.Find(&page2Users, bson.M{}, "users",
		ref.WithSkip(10), // Skip first 10 users
		ref.WithLimit(5), // Get next 5 users
		ref.WithSort(bson.D{{Key: "name", Value: 1}}), // Sort alphabetically
	)
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d users on page 2 (skipped 10, limit 5)\n", len(page2Users))
	}
}

// findWithProjectionExample demonstrates Find with field projection
func findWithProjectionExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 5: Find with Projection ===")
	var userEmails []bson.M
	err := mongoManager.Find(&userEmails, bson.M{"active": true}, "users",
		ref.WithProjection(bson.D{
			{Key: "name", Value: 1},  // Include name
			{Key: "email", Value: 1}, // Include email
			{Key: "_id", Value: 0},   // Exclude _id
		}),
		ref.WithLimit(3),
	)
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d active users (name and email only)\n", len(userEmails))
		for i, user := range userEmails {
			fmt.Printf("  %d. %s - %s\n", i+1, user["name"], user["email"])
		}
	}
}

// complexFindExample demonstrates Find with multiple options combined
func complexFindExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 6: Complex Find with All Options ===")
	var complexResult []bson.M
	err := mongoManager.Find(&complexResult,
		// Filter: adults between 18-65
		bson.M{
			"age": bson.M{
				"$gte": 18,
				"$lt":  65,
			},
			"active": true,
		},
		"users",
		// Sort by age, then by name
		ref.WithSort(bson.D{
			{Key: "age", Value: 1},
			{Key: "name", Value: 1},
		}),
		// Pagination: skip 0, limit 10
		ref.WithSkip(0),
		ref.WithLimit(10),
		// Only get specific fields
		ref.WithProjection(bson.D{
			{Key: "name", Value: 1},
			{Key: "age", Value: 1},
			{Key: "email", Value: 1},
			{Key: "city", Value: 1},
			{Key: "_id", Value: 0},
		}),
	)
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d active adults (18-65) with complete info\n", len(complexResult))
		for i, user := range complexResult {
			fmt.Printf("  %d. %s (age %v) - %s, %s\n",
				i+1, user["name"], user["age"], user["email"], user["city"])
		}
	}
}

// findWithMultipleSortExample demonstrates Find with multiple sort fields
func findWithMultipleSortExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 7: Find with Multiple Sort Fields ===")
	var multiSortUsers []bson.M
	err := mongoManager.Find(&multiSortUsers, bson.M{}, "users",
		ref.WithSort(bson.D{
			{Key: "department", Value: 1}, // First sort by department (ascending)
			{Key: "salary", Value: -1},    // Then by salary (descending)
			{Key: "name", Value: 1},       // Finally by name (ascending)
		}),
		ref.WithLimit(5),
	)
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d users sorted by department, salary (desc), then name\n", len(multiSortUsers))
	}
}

// findWithSpecificConnectionExample demonstrates Find using specific connection
func findWithSpecificConnectionExample(mongoManager db.IMongoLib) {
	fmt.Println("\n=== Example 8: Using Specific Connection ===")
	var analyticsData []bson.M
	err := mongoManager.Find(&analyticsData,
		bson.M{"type": "daily"},
		"metrics",
		// Note: Using default connection for this example
		ref.WithSort(bson.D{{Key: "date", Value: -1}}),
		ref.WithLimit(7), // Last 7 days
	)
	if err != nil {
		log.Printf(printMessage, errMessage, err)
	} else {
		fmt.Printf("Found %d daily analytics records\n", len(analyticsData))
	}
}

// FindExamples demonstrates all Find operation options by calling individual example functions
func FindExamples() {
	// Initialize MongoDB manager
	mongoManager := db.NewMongo()
	defer mongoManager.Close()

	// Run all examples
	basicFindExample(mongoManager)
	findWithLimitExample(mongoManager)
	findWithSortExample(mongoManager)
	findWithPaginationExample(mongoManager)
	findWithProjectionExample(mongoManager)
	complexFindExample(mongoManager)
	findWithMultipleSortExample(mongoManager)
	findWithSpecificConnectionExample(mongoManager)
}
