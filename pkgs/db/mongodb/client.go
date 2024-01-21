package mongodb

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"latipe-transaction-service/config"
	"time"
)

type MongoClient struct {
	client  *mongo.Client
	rootCtx context.Context
}

// Open - creates a new Mongo
func OpenMongoDBConnection(cfg *config.Config) (*MongoClient, error) {
	ctx := context.Background()

	if cfg.DB.Mongodb.ConnectTimeout == 0 {
		cfg.DB.Mongodb.ConnectTimeout = 15 * time.Second
	}

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(cfg.DB.Mongodb.Connection).
			SetConnectTimeout(cfg.DB.Mongodb.ConnectTimeout).
			SetMaxPoolSize(cfg.DB.Mongodb.MaxPoolSize))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoClient{client: client, rootCtx: ctx}, nil

}

func (m *MongoClient) GetDB() *mongo.Database {
	db := m.client.Database("latipe_delivery_db")
	return db
}

// Disconnect - used mainly in testing to avoid capping out the concurrent connections on MongoDB
func (m *MongoClient) Disconnect() {
	err := m.client.Disconnect(m.rootCtx)
	if err != nil {
		log.Fatalf("disconnecting from mongodb: %v", err)
	}
}

// Ping sends a ping command to verify that the client can connect to the deployment.
func (m *MongoClient) Ping() error {
	return m.client.Ping(m.rootCtx, readpref.Primary())
}
