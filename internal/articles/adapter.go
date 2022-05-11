package articles

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/multierr"
	"pokergo/pkg/id"
	"pokergo/pkg/pointers"
)

type Adapter interface {
	// Save saves articles (ignores duplicates)
	Save(ctx context.Context, arts []Article) ([]interface{}, error)
	// GetAll returns all saved articles
	GetAll(ctx context.Context) ([]shaArticle, error)
	// GetNext returns n documents
	GetNext(ctx context.Context, lastDocID id.ID, no int) ([]shaArticle, error)
}

type mongoAdapter struct {
	coll *mongo.Collection
}

func NewMongoAdapter(coll *mongo.Collection) *mongoAdapter {
	return &mongoAdapter{coll: coll}
}

func (m *mongoAdapter) EnsureIndexes(ctx context.Context) error {
	unique := options.IndexOptions{
		Unique: pointers.Pointer(true),
	}
	userIDIdx := mongo.IndexModel{
		Keys: bson.M{
			"hash_code": 1,
		},
		Options: &unique,
	}

	_, err := m.coll.Indexes().CreateOne(ctx, userIDIdx)
	if err != nil {
		return fmt.Errorf("cannot create unique name:1 index: %w", err)
	}

	return nil
}

func (m *mongoAdapter) Save(ctx context.Context, arts []Article) ([]interface{}, error) {
	if len(arts) == 0 {
		return nil, nil
	}

	var mongoArts []any // must be generic type
	for idx := range arts {
		mArt, err := convertToSHA(arts[idx])
		if err != nil {
			return nil, fmt.Errorf("cannot convert article to shaArticle: %w", err)
		}
		mongoArts = append(mongoArts, mArt)
	}

	opts := options.InsertMany().SetOrdered(false)
	res, insertErr := m.coll.InsertMany(ctx, mongoArts, opts)

	// filter-out E11000 (duplicate key error collection) as not-an-error
	var realErrors []error
	if insertErr != nil {
		bulkWrtErr, ok := insertErr.(mongo.BulkWriteException) // nolint:errorlint // this won't be a wrapped error
		if !ok {
			return nil, fmt.Errorf("mongo write err: %w", insertErr)
		}

		if wrErr := bulkWrtErr.WriteConcernError; wrErr != nil {
			realErrors = append(realErrors, wrErr)
		}

		for _, e := range bulkWrtErr.WriteErrors {
			if e.Code != 11000 { // duplicate key error
				realErrors = append(realErrors, e)
			}
		}
	}

	return res.InsertedIDs, multierr.Combine(realErrors...)
}

func (m *mongoAdapter) GetAll(ctx context.Context) ([]shaArticle, error) {
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("cannot perform the query: %w", err)
	}

	var articles []shaArticle
	if err := cur.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("cannot bind the data: %w", err)
	}

	return articles, nil
}

func (m *mongoAdapter) GetNext(ctx context.Context, lastDocID id.ID, no int) ([]shaArticle, error) {
	filter := bson.M{
		"_id": bson.M{
			"$gt": lastDocID,
		},
	}
	opts := &options.FindOptions{
		Limit: pointers.Pointer(int64(no)),
	}

	cur, err := m.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot get articles: %w", err)
	}

	var articles []shaArticle
	if err := cur.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("cannot bind articles: %w", err)
	}

	return articles, nil
}
