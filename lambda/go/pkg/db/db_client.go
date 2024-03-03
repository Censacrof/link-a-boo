package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DbClient struct {
	DdbClient *dynamodb.Client
}

var dbClient *DbClient = nil

func newDbClient(ctx context.Context) (*DbClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &DbClient{
		DdbClient: dynamodb.NewFromConfig(cfg),
	}, nil
}

func GetDbClient(ctx context.Context) (*DbClient, error) {
	if dbClient == nil {
		c, err := newDbClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("Can't create new Client: %w", err)
		}

		dbClient = c
	}

	return dbClient, nil
}
