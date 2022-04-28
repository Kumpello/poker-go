package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Collections struct {
	Users *mongo.Collection
}

func NewMongo(ctx context.Context, uri, authDB, user, pass, db string) (*Collections, error) {
	opts := options.Client().ApplyURI(uri)
	opts.SetAuth(options.Credential{
		AuthSource: authDB,
		Username:   user,
		Password:   pass,
	})
	opts.SetServerSelectionTimeout(time.Duration(30) * time.Second)

	cl, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("cannot create mongo client: %w", err)
	}

	if err := cl.Connect(ctx); err != nil {
		return nil, fmt.Errorf("cannot connect to the db: %w", err)
	}

	if err := cl.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping db error: %w", err)
	}

	appDB := cl.Database(db)

	return &Collections{
		Users: appDB.Collection("users"),
	}, nil
}
