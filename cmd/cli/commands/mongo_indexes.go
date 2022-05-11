package commands

import (
	"pokergo/internal/articles"
	"pokergo/internal/game"
	"pokergo/internal/org"
	"pokergo/internal/users"
)

func (c *commandApp) mongoIndexes() {
	usersAdapter := users.NewMongoAdapter(c.mongoColls.Users, c.logger)
	if err := usersAdapter.EnsureIndexes(c.Context()); err != nil {
		c.logger.Fatalf("cannot create indexes on users collection: %s", err.Error())
	}
	orgAdapter := org.NewMongoAdapter(c.mongoColls.Org, c.timer)
	if err := orgAdapter.EnsureIndexes(c.Context()); err != nil {
		c.logger.Fatalf("cannot create indexes on organizations collection: %s", err.Error())
	}
	gameAdapter := game.NewMongoAdapter(c.mongoColls.Games, c.timer)
	if err := gameAdapter.EnsureIndexes(c.Context()); err != nil {
		c.logger.Fatalf("cannot create indexes on games collection: %s", err.Error())
	}
	artsAdapter := articles.NewMongoAdapter(c.mongoColls.Arts)
	if err := artsAdapter.EnsureIndexes(c.Context()); err != nil {
		c.logger.Fatalf("cannot create indexes on articles collection: %s", err.Error())
	}
}
