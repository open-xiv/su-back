package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
	"time"
)

func ConnectDB() *mongo.Client {
	// options
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		zap.L().Fatal("failed to create mongo client", zap.Error(err))
		return nil
	}

	// connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ping
	err = client.Ping(ctx, nil)
	if err != nil {
		zap.L().Fatal("failed to ping mongo", zap.Error(err))
		return nil
	}
	zap.L().Debug("connected to mongo")
	return client
}
