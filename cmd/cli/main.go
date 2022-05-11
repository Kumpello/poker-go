package main

import (
	"context"

	"pokergo/cmd/cli/commands"
	"pokergo/internal/articles"
	"pokergo/internal/mongo"
	"pokergo/pkg/env"
	"pokergo/pkg/logger"
	"pokergo/pkg/timer"
)

func main() {
	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()

	mongoURI := env.Env("MONGO_URI", "mongodb://localhost:27017")
	mongoAuthDB := env.Env("MONGO_AUTH_DB", "admin")
	mongoUser := env.Env("MONGO_USER", "root")
	mongoPassword := env.Env("MONGO_PASSWORD", "password123")
	mongoDB := env.Env("MONGO_DB", "pokergo")
	mongoCollections, err := mongo.NewMongo(appCtx, mongoURI, mongoAuthDB, mongoUser, mongoPassword, mongoDB)
	if err != nil {
		log.Fatalf("cannot init mongo: %s", err.Error())
	}

	artsAdapter := articles.NewMongoAdapter(mongoCollections.Arts)

	cmd := commands.NewCommandApp(log, utcTimer, artsAdapter, mongoCollections)
	if err := cmd.ExecuteContext(appCtx); err != nil {
		log.Fatal(err)
	}
}
