# MongoDB Helper for Go (MongoDB Driver v2)

This package provides a reusable MongoDB helper with connection pooling and support for multiple MongoDB connections, compatible with MongoDB Go driver v2.

## Features

- Connection pooling with configurable pool size
- Support for multiple MongoDB connections with a simplified interface
- Thread-safe client management
- Configurable timeouts and connection parameters
- Functional options pattern for flexible configuration
- Easy access to databases and collections

## Usage

### Basic Usage

```go
// Initialize MongoDB manager
mongoManager := db.NewMongo()
defer mongoManager.Close()

// Simple operations with default connection
var user bson.M
err := mongoManager.FindOne("users", bson.M{"email": "john@example.com"}, &user)

// Insert a document
result, err := mongoManager.InsertOne("users", bson.M{"name": "John", "email": "john@example.com"})
```

### Multiple Connections

```go
// Add additional connections
err := mongoManager.AddConnection("analytics", db.MongoConfig{
    uri:            "mongodb://analytics-db:27017",
    maxPoolSize:    30,
    minPoolSize:    5,
    maxIdleTime:    10 * time.Minute,
    connectTimeout: 15 * time.Second,
    socketTimeout:  45 * time.Second,
})

// Use specific connection with functional options
var analyticsData bson.M
err = mongoManager.FindOne(
    &analyticsData,
    bson.M{"type": "daily"}, 
    "metrics",
    // Note: For multiple connections, you would configure them during setup
)
```

### Find Multiple Documents

```go
// Basic find
var users []bson.M
err := mongoManager.Find(&users, bson.M{"active": true}, "users")

// Find with limit and sort
var recentUsers []bson.M
err = mongoManager.Find(&recentUsers, bson.M{"active": true}, "users",
    db.WithLimit(10),
    db.WithSort(bson.D{{Key: "created", Value: -1}}),
)

// Find with pagination
var pagedUsers []bson.M
err = mongoManager.Find(&pagedUsers, bson.M{}, "users",
    db.WithSkip(20),      // Skip first 20 documents
    db.WithLimit(10),     // Return next 10 documents
    db.WithSort(bson.D{{Key: "name", Value: 1}}),
)

// Find with projection (only specific fields)
var userNames []bson.M
err = mongoManager.Find(&userNames, bson.M{"active": true}, "users",
    db.WithProjection(bson.D{
        {Key: "name", Value: 1},
        {Key: "email", Value: 1},
        {Key: "_id", Value: 0},
    }),
    db.WithLimit(5),
)

// Complex query with all options
var complexResult []bson.M
err = mongoManager.Find(&complexResult, 
    bson.M{"age": bson.M{"$gte": 18, "$lt": 65}}, 
    "users",
    db.WithSort(bson.D{{Key: "age", Value: 1}, {Key: "name", Value: 1}}),
    db.WithSkip(0),
    db.WithLimit(25),
    db.WithProjection(bson.D{
        {Key: "name", Value: 1},
        {Key: "age", Value: 1},
        {Key: "email", Value: 1},
        {Key: "_id", Value: 0},
    }),
)
```

### Available Find Options

- **WithLimit(n)**: Limit the number of documents returned
- **WithSkip(n)**: Skip the first n documents (useful for pagination)
- **WithSort(sort)**: Sort documents by specified fields
- **WithProjection(fields)**: Include/exclude specific fields from results
- **WithConnection(name)**: Use a specific database connection
- **WithDatabase(dbName)**: Use a specific database

### Sort Examples

```go
// Sort by single field (ascending)
db.WithSort(bson.D{{Key: "name", Value: 1}})

// Sort by single field (descending)
db.WithSort(bson.D{{Key: "created", Value: -1}})

// Sort by multiple fields
db.WithSort(bson.D{
    {Key: "age", Value: 1},     // First by age ascending
    {Key: "name", Value: 1},   // Then by name ascending
})
```

### Projection Examples

```go
// Include only specific fields
db.WithProjection(bson.D{
    {Key: "name", Value: 1},
    {Key: "email", Value: 1},
})

// Exclude specific fields
db.WithProjection(bson.D{
    {Key: "password", Value: 0},
    {Key: "internal_notes", Value: 0},
})

// Include fields and exclude _id
db.WithProjection(bson.D{
    {Key: "name", Value: 1},
    {Key: "email", Value: 1},
    {Key: "_id", Value: 0},
})
```

## Configuration

The default connection uses environment variables:

- `MONGO_URI`: MongoDB connection string (required)
- `MONGO_DB_NAME`: Default database name (required for operations)
- `MONGO_MAX_POOL_SIZE`: Maximum connection pool size (default: 20)
- `MONGO_MIN_POOL_SIZE`: Minimum connection pool size (default: 5)
- `MONGO_MAX_IDLE_TIME`: Maximum idle time in minutes (default: 5)

## Best Practices

1. **Initialize once**: Create a single instance of `MongoLib` at application startup
2. **Close properly**: Always call `Close()` when shutting down your application
3. **Configure pool size**: Adjust pool sizes based on your application's needs
4. **Use timeouts**: Set appropriate timeouts for your operations
5. **Handle errors**: Always check for errors when performing database operations
6. **Use default connection**: For most operations, you don't need to specify a connection name
7. **Use functional options**: The functional options pattern makes the API more flexible and readable

See `mongo_example.go` for a complete example.
