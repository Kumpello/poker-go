package org

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pokergo/pkg/id"
	"pokergo/pkg/pointers"
	"pokergo/pkg/timer"
	"time"
)

type Org struct {
	ID        id.ID     `bson:"_id"`
	Name      string    `bson:"name"`
	Admin     id.ID     `bson:"admin"`
	Members   []id.ID   `bson:"members"`
	CreatedAt time.Time `bson:"created_at"`
}

func (o Org) IsMember(id id.ID) bool {
	if id == o.Admin {
		return true
	}
	for idx := range o.Members {
		if o.Members[idx] == id {
			return true
		}
	}
	return false
}

type Adapter interface {
	// GetOrgByID returns org by id
	GetOrgByID(ctx context.Context, id id.ID) (Org, error)
	// GetOrgByName returns org by its name
	GetOrgByName(ctx context.Context, name string) (Org, error)
	// CreateOrg creates a new Org
	CreateOrg(ctx context.Context, admin id.ID, orgName string) (Org, error)
	// AddToOrg adds a new member to the organization
	AddToOrg(ctx context.Context, orgID id.ID, who id.ID) error
	// ListUserOrg list all organizations where user belongs to
	ListUserOrg(ctx context.Context, userID id.ID) ([]Org, error)
}

type mongoAdapter struct {
	coll  *mongo.Collection
	timer timer.Timer
}

var ErrOrgNotExists = mongo.ErrNoDocuments

func NewMongoAdapter(coll *mongo.Collection, timer timer.Timer) *mongoAdapter {
	return &mongoAdapter{coll: coll, timer: timer}
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

func (m *mongoAdapter) GetOrgByID(ctx context.Context, id id.ID) (Org, error) {
	filter := bson.M{
		"_id": id,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Org{}, ErrOrgNotExists
		}
		return Org{}, fmt.Errorf("cannot find org")
	}

	var org Org
	if err := res.Decode(&org); err != nil {
		return Org{}, fmt.Errorf("cannot decode result: %w", err)
	}

	return org, nil
}

func (m *mongoAdapter) GetOrgByName(ctx context.Context, name string) (Org, error) {
	filter := bson.M{
		"name": name,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Org{}, ErrOrgNotExists
		}
		return Org{}, fmt.Errorf("cannot find org")
	}

	var org Org
	if err := res.Decode(&org); err != nil {
		return Org{}, fmt.Errorf("cannot decode result: %w", err)
	}

	return org, nil
}

func (m *mongoAdapter) CreateOrg(ctx context.Context, admin id.ID, orgName string) (Org, error) {
	newOrg := Org{
		ID:        id.NewID(),
		Name:      orgName,
		Admin:     admin,
		Members:   []id.ID{admin},
		CreatedAt: m.timer.Now(),
	}

	_, err := m.coll.InsertOne(ctx, newOrg)
	if err != nil {
		return Org{}, fmt.Errorf("cannot insert a new organization: %w", err)
	}

	return newOrg, nil
}

func (m *mongoAdapter) AddToOrg(ctx context.Context, orgID id.ID, who id.ID) error {
	find := bson.M{
		"_id": orgID,
	}
	update := bson.M{
		"$push": bson.M{
			"members": who,
		},
	}

	res, err := m.coll.UpdateOne(ctx, find, update)
	if err != nil {
		return fmt.Errorf("cannot update members: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find the organization")
	}

	return nil
}

func (m *mongoAdapter) ListUserOrg(ctx context.Context, userID id.ID) ([]Org, error) {
	find := bson.M{
		"members": bson.M{
			"$in": []any{userID},
		},
	}

	cur, err := m.coll.Find(ctx, find)
	if err != nil {
		return nil, fmt.Errorf("cannot perform find query: %w", err)
	}

	var result []Org
	if err := cur.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("cannot bind query result: %w", err)
	}

	return result, nil
}
