package users

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pokergo/pkg/id"
	"pokergo/pkg/logger"
	"pokergo/pkg/pointers"
)

type User struct {
	// ID is internal ID
	ID id.ID `bson:"_id"` // nolint:tagliatelle // mongo-id
	// Username is a unique name of player (nick), used for login as well
	Username string `bson:"name"`
	// Email is user email
	Email string `bson:"email"`
	// Password is an encrypted password
	Password string `bson:"password"`
	// Token is a jwt-token
	Token string `bson:"token"`
	// RefreshToken is jwt-refresh-token
	RefreshToken string `bson:"refresh_token"`
	// CreatedAt tells when the user was created
	CreatedAt time.Time `bson:"created_at"`
	// UpdatedAt tells when the last update was done
	UpdatedAt time.Time `bson:"updated_at"`
}

type Adapter interface {
	NewUser(ctx context.Context, user User) (User, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	GetUserByID(ctx context.Context, id id.ID) (User, error)

	UserDetails(ctx context.Context, ids []id.ID) (map[id.ID]User, error)

	// UpdateTokens update user tokens (if they are not nil)
	UpdateTokens(ctx context.Context, userID id.ID, token, refreshedToken *string) error
}

var ErrUserNotExists = mongo.ErrNoDocuments

type mongoAdapter struct {
	coll   *mongo.Collection
	logger logger.Logger
}

func NewMongoAdapter(coll *mongo.Collection, logger logger.Logger) *mongoAdapter {
	return &mongoAdapter{coll, logger}
}

func (m *mongoAdapter) EnsureIndexes(ctx context.Context) error {
	unique := options.IndexOptions{
		Unique: pointers.Pointer(true),
	}
	userIDIdx := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: &unique,
	}

	_, err := m.coll.Indexes().CreateOne(ctx, userIDIdx)
	if err != nil {
		return fmt.Errorf("cannot create unique name:1 index: %w", err)
	}

	return nil
}

func (m *mongoAdapter) NewUser(ctx context.Context, user User) (User, error) {
	user.ID = id.NewID()
	_, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("cannot create a new user: %w", err)
	}

	return user, nil
}

func (m *mongoAdapter) GetUserByName(ctx context.Context, name string) (User, error) {
	filter := bson.M{
		"name": name,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, fmt.Errorf("cannot perform query: %w", err)
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, fmt.Errorf("cannot decode query result: %w", err)
	}

	return user, nil
}

func (m *mongoAdapter) GetUserByID(ctx context.Context, id id.ID) (User, error) {
	filter := bson.M{
		"_id": id,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, fmt.Errorf("cannot perform query: %w", err)
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, fmt.Errorf("cannot decode query result: %w", err)
	}

	return user, nil
}

func (m *mongoAdapter) UserDetails(ctx context.Context, ids []id.ID) (map[id.ID]User, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	c, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("cannot perform find query: %w", err)
	}

	var users []User
	if err := c.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("cannot decode query result: %w", err)
	}

	res := make(map[id.ID]User, len(users))
	for _, v := range users {
		res[v.ID] = v
	}
	return res, nil
}

func (m *mongoAdapter) UpdateTokens(ctx context.Context, userID id.ID, token, refreshedToken *string) error {
	tokens := bson.D{}
	if token != nil {
		tokens = append(tokens, bson.E{Key: "token", Value: *token})
	}
	if refreshedToken != nil {
		tokens = append(tokens, bson.E{Key: "refresh_token", Value: *refreshedToken})
	}

	if len(tokens) == 0 {
		// nothing to update
		return nil
	}

	filter := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$set": tokens,
	}

	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("cannot update tokens: %w", err)
	}

	return nil
}

var _ Adapter = (*mongoAdapter)(nil)
