# MongoDB Upsert Functionality

This document explains how to use the upsert functionality in the MongoDB library.

## What is Upsert?

Upsert is a database operation that combines "update" and "insert". If a document matching the filter exists, it updates the document. If no document matches the filter, it creates a new document.

## Usage

### Basic Upsert with UpdateOneSet

```go
import (
    "github.com/ranggadablues/gosok/db"
    "github.com/ranggadablues/gosok/db/ref"
    "go.mongodb.org/mongo-driver/v2/bson"
)

// Initialize MongoDB manager
mongoManager := db.NewMongo()
defer mongoManager.Close()

// Update or create a user
updateData := bson.M{
    "name":    "John Doe",
    "email":   "john@example.com",
    "age":     30,
    "status":  "active",
}

err := mongoManager.UpdateOneSet(
    "users",
    bson.M{"email": "john@example.com"}, // filter
    updateData,                          // update data
    ref.WithUpsert(true),                // enable upsert
)
```

### Upsert with Pipeline Updates

```go
// Using pipeline format for more complex operations
pipelineUpdate := bson.M{
    "$set": bson.M{
        "name": "Jane Smith",
        "age":  28,
    },
    "$inc": bson.M{
        "login_count": 1,
    },
}

err := mongoManager.UpdateOneSetPipeline(
    "users",
    bson.M{"email": "jane@example.com"},
    pipelineUpdate,
    ref.WithUpsert(true),
)
```

### Upsert with UpdateManySet

```go
// Update multiple documents with upsert
err := mongoManager.UpdateManySet(
    "users",
    bson.M{"status": "inactive"},
    bson.M{
        "status": "active",
        "updated": time.Now(),
    },
    ref.WithUpsert(true),
)
```

## Key Points

1. **Backward Compatibility**: All existing code continues to work without changes. The upsert option is optional.

2. **Filter-based**: The upsert operation uses the filter to determine if a document exists. If a document matches the filter, it gets updated. If no document matches, a new one is created.

3. **Update Data**: When creating a new document (upsert), the update data is used to populate the new document.

4. **Options**: Use `ref.WithUpsert(true)` to enable upsert functionality.

## Examples

See the following files for complete examples:
- `update_examples.go` - Comprehensive examples of all upsert scenarios
- `test_upsert.go` - Simple test to verify upsert functionality

## Common Use Cases

1. **User Management**: Create or update user profiles
2. **Configuration**: Set default configurations that can be updated later
3. **Counters**: Initialize counters that can be incremented
4. **Caching**: Store computed values that can be refreshed

## Error Handling

Always check for errors when performing upsert operations:

```go
err := mongoManager.UpdateOneSet(
    "users",
    filter,
    updateData,
    ref.WithUpsert(true),
)
if err != nil {
    log.Printf("Upsert failed: %v", err)
    return
}
```

## Performance Considerations

- Upsert operations are slightly slower than regular updates
- Use appropriate indexes on filter fields for better performance
- Consider the impact on write operations in high-traffic scenarios
