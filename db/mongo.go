package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ranggadablues/gosok/db/ref"
	"github.com/ranggadablues/gosok/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
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
	Debug() *MongoLib

	// Database operations
	FindOne(output, filter any, collName string, opts ...ref.FindOption) error
	Find(output, filter any, collName string, opts ...ref.FindOption) error
	InsertOne(collName string, document any) (any, error)
	InsertMany(collName string, documents []any) ([]any, error)
	DeleteOne(collName string, filter any) error
	DeleteMany(collName string, filter any) error
	updateOne(collName string, filter any, update any, opts ...ref.UpdateOption) error
	UpdateOneSet(collName string, filter any, update any, opts ...ref.UpdateOption) error
	UpdateOneSetPipeline(collName string, filter any, update any, opts ...ref.UpdateOption) error
	updateMany(collName string, filter any, update any, opts ...ref.UpdateOption) error
	UpdateManySet(collName string, filter any, update any, opts ...ref.UpdateOption) error
	UpdateManySetPipeline(collName string, filter any, update any, opts ...ref.UpdateOption) error
	Aggregate(output, pipeline any, collName string) error
}

// MongoLib manages a single MongoDB connection
type MongoLib struct {
	uri        string
	client     *mongo.Client
	database   *mongo.Database
	ctx        context.Context
	logger     func() logger.ILogLevel
	isdebug    bool
	isconninfo bool
}

// NewMongo creates a new MongoDB connection
// if args[0] is true, set isconninfo to true
func NewMongo(args ...bool) IMongoLib {
	m := &MongoLib{
		ctx:        context.Background(),
		logger:     logger.NewLogger,
		isdebug:    false,
		isconninfo: false,
	}

	if len(args) > 0 {
		m.isconninfo = args[0]
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

	if m.isconninfo {
		clientOpts.SetPoolMonitor(m.setPoolMonitor())
		clientOpts.SetMonitor(m.setMonitor())
	}

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
	m.logger().UTC().LogInfoLevel("msg", "MongoDB connected successfully")

	return nil
}

func (m *MongoLib) setPoolMonitor() *event.PoolMonitor {
	// Monitor pool connections
	poolMonitor := &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			switch evt.Type {
			case event.ConnectionCreated:
				print := fmt.Sprintf("[POOL] Connection created: id=%d, address=%s", evt.ConnectionID, evt.Address)
				m.logger().LogInfoLevel("msg", print)
			case event.ConnectionReady:
				print := fmt.Sprintf("[POOL] Connection ready: id=%d", evt.ConnectionID)
				m.logger().LogInfoLevel("msg", print)
			case event.ConnectionClosed:
				print := fmt.Sprintf("[POOL] Connection closed: id=%d, reason=%s", evt.ConnectionID, evt.Reason)
				m.logger().LogInfoLevel("msg", print)
			case event.ConnectionCheckedOut:
				print := fmt.Sprintf("[POOL] Connection checked out: id=%d", evt.ConnectionID)
				m.logger().LogInfoLevel("msg", print)
			case event.ConnectionCheckedIn:
				print := fmt.Sprintf("[POOL] Connection checked in: id=%d", evt.ConnectionID)
				m.logger().LogInfoLevel("msg", print)
			}
		},
	}

	return poolMonitor
}

func (m *MongoLib) setMonitor() *event.CommandMonitor {
	// Monitor commands (queries)
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			print := fmt.Sprintf("[QUERY] %s on %s cmd=%v", evt.CommandName, evt.DatabaseName, evt.Command)
			m.logger().LogInfoLevel("msg", print)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			print := fmt.Sprintf("[QUERY] Done %s (%dms)", evt.CommandName, evt.Duration.Milliseconds())
			m.logger().LogInfoLevel("msg", print)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			print := fmt.Sprintf("[QUERY] FAIL %s (%v)", evt.CommandName, evt.Failure)
			m.logger().LogInfoLevel("msg", print)
		},
	}

	return cmdMonitor
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

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("FindOne")
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

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("FindMany")
	}

	return cursor.All(m.ctx, output)
}

// InsertOne inserts a single document into the specified collection
func (m *MongoLib) InsertOne(collName string, document any) (any, error) {
	if err := m.ensureConnection(); err != nil {
		return bson.NilObjectID, err
	}
	collection := m.GetCollection(collName)
	result, err := collection.InsertOne(m.ctx, document)
	if err != nil {
		return bson.NilObjectID, err
	}
	if !result.Acknowledged {
		return bson.NilObjectID, errors.New("insert not acknowledged")
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("InsertOne")
	}

	return result.InsertedID, nil
}

// InsertMany inserts multiple documents into the specified collection
func (m *MongoLib) InsertMany(collName string, documents []any) ([]any, error) {
	if err := m.ensureConnection(); err != nil {
		return nil, err
	}
	collection := m.GetCollection(collName)
	result, err := collection.InsertMany(m.ctx, documents)
	if err != nil {
		return nil, err
	}
	if !result.Acknowledged {
		return nil, errors.New("insert not acknowledged")
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("InsertMany")
	}

	return result.InsertedIDs, nil
}

// DeleteOne deletes a single document from the specified collection
func (m *MongoLib) DeleteOne(collName string, filter any) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}
	collection := m.GetCollection(collName)
	result, err := collection.DeleteOne(m.ctx, filter)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return errors.New("delete not acknowledged")
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("DeleteOne")
	}

	return nil
}

// DeleteMany deletes multiple documents from the specified collection
func (m *MongoLib) DeleteMany(collName string, filter any) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}
	collection := m.GetCollection(collName)
	result, err := collection.DeleteMany(m.ctx, filter)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return errors.New("delete not acknowledged")
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("DeleteMany")
	}

	return nil
}

// UpdateOneSet(collName string, filter any, update any, opts ...ref.UpdateOption) error
// e.g db.collectionName.update({_id: "123"}, {$set: {name: "John"}})
func (m *MongoLib) UpdateOneSet(collName string, filter any, update any, opts ...ref.UpdateOption) error {
	return m.updateOne(collName, filter, ref.UpdateSet(update), opts...)
}

// UpdateOneSetPipeline(collName string, filter any, update any, opts ...ref.UpdateOption) error
// e.g db.collectionName.update({_id: "123"}, [{$set: {name: "$otherfield"}}])
func (m *MongoLib) UpdateOneSetPipeline(collName string, filter any, update any, opts ...ref.UpdateOption) error {
	return m.updateOne(collName, filter, ref.UpdateSetPipeline(update), opts...)
}

// UpdateOne updates a single document in the specified collection
func (m *MongoLib) updateOne(collName string, filter any, update any, opts ...ref.UpdateOption) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}

	// Parse update options
	updateOpts := &ref.UpdateOptions{
		Upsert: nil,
	}

	// Apply options
	for _, opt := range opts {
		opt(updateOpts)
	}

	collection := m.GetCollection(collName)

	// Build MongoDB update options
	mongoOpts := options.UpdateOne()
	if updateOpts.Upsert != nil {
		mongoOpts.SetUpsert(*updateOpts.Upsert)
	}

	result, err := collection.UpdateOne(m.ctx, filter, update, mongoOpts)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return errors.New("update not acknowledged")
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("UpdateOne")
	}

	return nil
}

// UpdateManySet(collName string, filter any, update any, opts ...ref.UpdateOption) error
// e.g db.collectionName.updateMany({_id: "123"}, {$set: {name: "John"}})
func (m *MongoLib) UpdateManySet(collName string, filter any, update any, opts ...ref.UpdateOption) error {
	return m.updateMany(collName, filter, ref.UpdateSet(update), opts...)
}

// UpdateManySetPipeline(collName string, filter any, update any, opts ...ref.UpdateOption) error
// e.g db.collectionName.updateMany({_id: "123"}, [{$set: {name: "$otherfield"}}])
func (m *MongoLib) UpdateManySetPipeline(collName string, filter any, update any, opts ...ref.UpdateOption) error {
	return m.updateMany(collName, filter, ref.UpdateSetPipeline(update), opts...)
}

// UpdateMany updates multiple documents in the specified collection
func (m *MongoLib) updateMany(collName string, filter any, update any, opts ...ref.UpdateOption) error {
	if err := m.ensureConnection(); err != nil {
		return err
	}

	// Parse update options
	updateOpts := &ref.UpdateOptions{
		Upsert: nil,
	}

	// Apply options
	for _, opt := range opts {
		opt(updateOpts)
	}

	collection := m.GetCollection(collName)

	// Build MongoDB update options
	mongoOpts := options.UpdateMany()
	if updateOpts.Upsert != nil {
		mongoOpts.SetUpsert(*updateOpts.Upsert)
	}

	result, err := collection.UpdateMany(m.ctx, filter, update, mongoOpts)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return errors.New("update not acknowledged")
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("UpdateMany")
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

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("Aggregate")
	}

	return cursor.All(m.ctx, output)
}

// Count counts the number of documents in the specified collection
func (m *MongoLib) Count(collName string, filter any) (int64, error) {
	if err := m.ensureConnection(); err != nil {
		return 0, err
	}
	collection := m.GetCollection(collName)
	count, err := collection.CountDocuments(m.ctx, filter)
	if err != nil {
		return 0, err
	}

	if m.isdebug {
		m.logger().UTC().LogDebugLevelWithCaller("CountDocuments")
	}

	return count, nil
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
		m.logger().UTC().LogWarnLevel("msg", "Connection lost, attempting to reconnect:", err.Error())
		// Try to reconnect
		return m.connect()
	}

	return nil
}

func (m *MongoLib) Debug() *MongoLib {
	m.isdebug = true
	return m
}
