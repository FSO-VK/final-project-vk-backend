// Package mongo is an implementation of vidal storage for MongoDB.
package mongo

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/logcon"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const productCollection = "products"

// Storage implements vidal.Storage for a MongoDB.
type Storage struct {
	client   *mongo.Client
	database *mongo.Database
	config   *Config
	log      *logrus.Entry
}

// NewStorage creates a new MongoDB storage.
func NewStorage(config *Config, log *logrus.Entry) (*Storage, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s/%s",
		config.User,
		config.Password,
		net.JoinHostPort(config.Host, config.Port), // required if using IPv6
		config.Database,
	)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("MongoDB connect: %w", err)
	}

	s := &Storage{
		client:   client,
		database: client.Database(config.Database),
		config:   config,
		log:      log,
	}
	s.debugLog(context.Background(), "MongoDB storage DEBUG mode enabled")

	err = client.Ping(context.Background(), nil)
	if err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("MongoDB ping: %w", err)
	}
	s.debugLog(context.Background(), "MongoDB storage connected")

	err = s.initConfig()
	if err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("init config: %w", err)
	}
	s.debugLog(context.Background(), "MongoDB storage initialized")

	return s, nil
}

// SaveProduct saves a product in MongoDB.
func (s *Storage) SaveProduct(ctx context.Context, product *vidal.StorageModel) error {
	s.debugLog(ctx, "save product: %v", product)

	coll := s.database.Collection(productCollection)
	saved, err := coll.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}

	s.debugLog(ctx, "saved: %t", saved.Acknowledged)
	return nil
}

// GetProduct loads a product from MongoDB.
func (s *Storage) GetProduct(ctx context.Context, barCode string) (*vidal.StorageModel, error) {
	s.debugLog(ctx, "load product by bar code %s", barCode)
	coll := s.database.Collection(productCollection)

	var result vidal.StorageModel
	err := coll.FindOne(
		ctx,
		bson.D{{Key: "barCodes", Value: barCode}},
	).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, vidal.ErrStorageNoProduct
		}
		return nil, fmt.Errorf("find product decode: %w", err)
	}
	s.debugLog(ctx, "loaded: %v", result)

	return &result, nil
}

// Close gracefully closes MongoDB connection.
func (s *Storage) Close(ctx context.Context) error {
	if s.client != nil {
		return s.client.Disconnect(ctx)
	}
	return nil
}

// initConfig performs initial configuration of the database.
func (s *Storage) initConfig() error {
	coll := s.database.Collection(productCollection)

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "barCodes", Value: 1}},
		Options: nil,
	}
	idxName, err := coll.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return fmt.Errorf("create index %s: %w", idxName, err)
	}

	s.debugLog(context.Background(), "index")
	return nil
}

func (s *Storage) debugLog(ctx context.Context, format string, args ...any) {
	if !s.config.Log {
		return
	}
	log, ok := logcon.FromContext(ctx)
	if !ok {
		log = s.log
	}
	log.Debugf(format, args...)
}
