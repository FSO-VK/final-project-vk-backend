package mongo

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	client   *mongo.Client
	database *mongo.Database
	config   *Config
}

func NewStorage(config *Config) (*Storage, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("create mongo client: %w", err)
	}

	return &Storage{
		client:   client,
		database: client.Database(config.Database),
		config:   config,
	}, nil
}

func (s *Storage) SaveProduct(ctx context.Context, product *vidal.StorageModel) error {
	col := s.database.Collection("products")
	_, err := col.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

