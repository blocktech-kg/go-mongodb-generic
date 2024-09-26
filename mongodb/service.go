package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, dbConnectionUrl string, dbName string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(dbConnectionUrl)

	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		err = fmt.Errorf("failed to mongo.Connect: %s", err)
		return nil, err
	}
	err = dbClient.Ping(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to dbClient.Ping: %s", err)
		return nil, err
	}

	return dbClient.Database(dbName), nil

}
