package examples

import (
	"sync"

	"github.com/ranggadablues/gosok/db"
)

var (
	mongoInstance db.IMongoLib
	mongoOnce     sync.Once
)

// GetMongoInstance returns a singleton instance of MongoDB connection
// This ensures all services share the same connection pool
func GetMongoInstance() db.IMongoLib {
	mongoOnce.Do(func() {
		mongoInstance = db.NewMongo()
	})
	return mongoInstance
}

// CloseMongoInstance closes the singleton MongoDB connection
// Call this during application shutdown
func CloseMongoInstance() error {
	if mongoInstance != nil {
		return mongoInstance.Close()
	}
	return nil
}
