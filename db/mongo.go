package db

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ranggadablues/gosok/db/ref"
	"github.com/ranggadablues/gosok/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// IMongoLib defines the interface for MongoDB operations
type IMongoLib interface {
	Close() error
	GetClient() *mongo.Client
	GetCollection(collName string) *mongo.Collection
	GetDatabaseName() string

	// Database operations
	FindOne(output, filter any, collName string, opts ...ref.FindOption) error
	Find(output, filter any, collName string, opts ...ref.FindOption) error
	InsertOne(collName string, document any) (*mongo.InsertOneResult, error)
	InsertMany(collName string, documents []any) (*mongo.InsertManyResult, error)
	DeleteOne(collName string, filter any) (*mongo.DeleteResult, error)
	DeleteMany(collName string, filter any) (*mongo.DeleteResult, error)
	updateOne(collName string, filter any, update any) error
	UpdateOneSet(collName string, filter any, update any) error
	UpdateOneSetPipeline(collName string, filter any, update any) error
	updateMany(collName string, filter any, update any) error
	UpdateManySet(collName string, filter any, update any) error
	UpdateManySetPipeline(collName string, filter any, update any) error
	Aggregate(output, pipeline any, collName string) error
}

// MongoLib manages a single MongoDB connection
type MongoLib struct {
	uri      string
	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context
	logger   func() logger.ILogLevel
}

// NewMongo creates a new MongoDB connection
func NewMongo() IMongoLib {
	m := &MongoLib{
		ctx:    context.Background(),
		logger: logger.NewLogger,
	}

	// Connect to MongoDB
	err := m.connect()
	if err != nil {
		m.logger().LogErrorLevel("msg", "error connecting to MongoDB:", err.Error())
		return nil
	}

	return m
}

// connect establishes a connection to MongoDB
func (m *MongoLib) connect() error {
	// Get MongoDB URI from environment
	m.uri = os.Getenv("MONGO_URI")
	if m.uri == "" {
		return errors.New("MONGO_URI environment variable is required")
	}

	// Get database name from environment
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		return errors.New("MONGO_DB_NAME environment variable is required")
	}

	// Configure client options with basic settings
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().
		ApplyURI(m.uri).
		SetMaxPoolSize(20).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Minute).
		SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return err
	}

	// Verify connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	// Store client and database
	m.client = client
	m.database = client.Database(dbName)
	m.logger().LogInfoLevel("msg", "MongoDB connected successfully")

	return nil
}

// GetClient returns the MongoDB client
func (m *MongoLib) GetClient() *mongo.Client {
	return m.client
}

// GetCollection returns a MongoDB collection
func (m *MongoLib) GetCollection(collName string) *mongo.Collection {
	return m.database.Collection(collName)
}

// GetDatabase returns a MongoDB database
func (m *MongoLib) GetDatabaseName() string {
	return m.database.Name()
}

// Close disconnects the MongoDB client
func (m *MongoLib) Close() error {
	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		m.logger().LogErrorLevel("msg", "Failed to disconnect from MongoDB:", err.Error())
		return err
	}

	m.logger().LogInfoLevel("msg", "MongoDB disconnected successfully")
	return nil
}

// FindOne finds a single document in the specified collection
func (m *MongoLib) FindOne(output, filter any, collName string, opts ...ref.FindOption) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}

	// Parse find options
	findOpts := &ref.FindOptions{
		Limit:      nil,
		Skip:       nil,
		Sort:       nil,
		Projection: nil,
	}

	// Apply options
	for _, opt := range opts {
		opt(findOpts)
	}

	// Get collection
	collection := m.GetCollection(collName)

	// Build MongoDB find options
	mongoOpts := options.FindOne()
	if findOpts.Sort != nil {
		mongoOpts.SetSort(findOpts.Sort)
	}
	if findOpts.Projection != nil {
		mongoOpts.SetProjection(findOpts.Projection)
	}
	if findOpts.Skip != nil {
		mongoOpts.SetSkip(*findOpts.Skip)
	}

	// Execute FindOne with options
	err := collection.FindOne(m.ctx, filter, mongoOpts).Decode(output)
	if err != nil {
		return err
	}
	return nil
}

// Find finds multiple documents in the specified collection
func (m *MongoLib) Find(output, filter any, collName string, opts ...ref.FindOption) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}

	// Parse find options
	findOpts := &ref.FindOptions{
		Limit:      nil,
		Skip:       nil,
		Sort:       nil,
		Projection: nil,
	}

	// Apply options
	for _, opt := range opts {
		opt(findOpts)
	}

	// Get collection
	collection := m.GetCollection(collName)

	// Build MongoDB find options
	mongoOpts := options.Find()
	if findOpts.Limit != nil {
		mongoOpts.SetLimit(*findOpts.Limit)
	}
	if findOpts.Skip != nil {
		mongoOpts.SetSkip(*findOpts.Skip)
	}
	if findOpts.Sort != nil {
		mongoOpts.SetSort(findOpts.Sort)
	}
	if findOpts.Projection != nil {
		mongoOpts.SetProjection(findOpts.Projection)
	}

	// Execute find with options
	cursor, err := collection.Find(m.ctx, filter, mongoOpts)
	if err != nil {
		return err
	}
	defer cursor.Close(m.ctx)

	return cursor.All(m.ctx, output)
}

// InsertOne inserts a single document into the specified collection
func (m *MongoLib) InsertOne(collName string, document any) (*mongo.InsertOneResult, error) {
	if err := m.ensureConnection(); err != nil {
		return nil, err
	}
	collection := m.GetCollection(collName)
	return collection.InsertOne(m.ctx, document)
}

// InsertMany inserts multiple documents into the specified collection
func (m *MongoLib) InsertMany(collName string, documents []any) (*mongo.InsertManyResult, error) {
	if err := m.ensureConnection(); err != nil {
		return nil, err
	}
	collection := m.GetCollection(collName)
	return collection.InsertMany(m.ctx, documents)
}

// DeleteOne deletes a single document from the specified collection
func (m *MongoLib) DeleteOne(collName string, filter any) (*mongo.DeleteResult, error) {
	if err := m.ensureConnection(); err != nil {
		return nil, err
	}
	collection := m.GetCollection(collName)
	return collection.DeleteOne(m.ctx, filter)
}

// DeleteMany deletes multiple documents from the specified collection
func (m *MongoLib) DeleteMany(collName string, filter any) (*mongo.DeleteResult, error) {
	if err := m.ensureConnection(); err != nil {
		return nil, err
	}
	collection := m.GetCollection(collName)
	return collection.DeleteMany(m.ctx, filter)
}

// UpdateOneSet(collName string, filter any, update any) error
// e.g db.collectionName.update({_id: "123"}, {$set: {name: "John"}})
func (m *MongoLib) UpdateOneSet(collName string, filter any, update any) error {
	return m.updateOne(collName, filter, ref.UpdateSet(update))
}

// UpdateOneSetPipeline(collName string, filter any, update any) error
// e.g db.collectionName.update({_id: "123"}, [{$set: {name: "$otherfield"}}])
func (m *MongoLib) UpdateOneSetPipeline(collName string, filter any, update any) error {
	return m.updateOne(collName, filter, ref.UpdateSetPipeline(update))
}

// UpdateOne updates a single document in the specified collection
func (m *MongoLib) updateOne(collName string, filter any, update any) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}
	collection := m.GetCollection(collName)
	result, err := collection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return errors.New("update not acknowledged")
	}

	return nil
}

// UpdateManySet(collName string, filter any, update any) error
// e.g db.collectionName.updateMany({_id: "123"}, {$set: {name: "John"}})
func (m *MongoLib) UpdateManySet(collName string, filter any, update any) error {
	return m.updateMany(collName, filter, ref.UpdateSet(update))
}

// UpdateManySetPipeline(collName string, filter any, update any) error
// e.g db.collectionName.updateMany({_id: "123"}, [{$set: {name: "$otherfield"}}])
func (m *MongoLib) UpdateManySetPipeline(collName string, filter any, update any) error {
	return m.updateMany(collName, filter, ref.UpdateSetPipeline(update))
}

// UpdateMany updates multiple documents in the specified collection
func (m *MongoLib) updateMany(collName string, filter any, update any) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}
	collection := m.GetCollection(collName)
	result, err := collection.UpdateMany(m.ctx, filter, update)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return errors.New("update not acknowledged")
	}
	return nil
}

// Aggregate aggregates documents from the specified collection
func (m *MongoLib) Aggregate(output, pipeline any, collName string) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}
	collection := m.GetCollection(collName)
	cursor, err := collection.Aggregate(m.ctx, pipeline)
	if err != nil {
		return err
	}

	return cursor.All(m.ctx, output)
}

// Count counts the number of documents in the specified collection
func (m *MongoLib) Count(collName string, filter any) (int64, error) {
	if err := m.ensureConnection(); err != nil {
		return 0, err
	}
	collection := m.GetCollection(collName)
	return collection.CountDocuments(m.ctx, filter)
}

// ensureConnection checks if connection is alive and reconnects if needed
func (m *MongoLib) ensureConnection() error {
	if m.client == nil {
		return m.connect()
	}

	// Ping to check if connection is still alive
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		m.logger().LogWarnLevel("msg", "Connection lost, attempting to reconnect:", err.Error())
		// Try to reconnect
		return m.connect()
	}

	return nil
}
