package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"pokergo/internal/game"
	"pokergo/internal/mongo"
	"pokergo/internal/org"
	"pokergo/internal/users"
	"pokergo/internal/webapi"
	authMux "pokergo/internal/webapi/auth"
	gameMux "pokergo/internal/webapi/game"
	newsMux "pokergo/internal/webapi/news"
	orgMux "pokergo/internal/webapi/org"
	"pokergo/pkg/env"
	"pokergo/pkg/jwt"
	"pokergo/pkg/logger"
	"pokergo/pkg/timer"
)

func main() {
	appCtx := context.Background()
	log := logger.NewLogger()
	utcTimer := timer.NewUTCTimer()

	// Mongo
	mongoURI := env.Env("MONGO_URI", "mongodb://localhost:27017")
	mongoAuthDB := env.Env("MONGO_AUTH_DB", "admin")
	mongoUser := env.Env("MONGO_USER", "root")
	mongoPassword := env.Env("MONGO_PASSWORD", "password123")
	mongoDB := env.Env("MONGO_DB", "pokergo")
	mongoCollections, err := mongo.NewMongo(appCtx, mongoURI, mongoAuthDB, mongoUser, mongoPassword, mongoDB)
	if err != nil {
		log.Fatalf("cannot init mongo: %s", err.Error())
	}
	usersAdapter := users.NewMongoAdapter(mongoCollections.Users, log)
	if err := usersAdapter.EnsureIndexes(appCtx); err != nil {
		log.Fatalf("cannot create indexes on users collection: %s", err.Error())
	}
	orgAdapter := org.NewMongoAdapter(mongoCollections.Org, utcTimer)
	if err := usersAdapter.EnsureIndexes(appCtx); err != nil {
		log.Fatalf("cannot create indexes on organizations collection: %s", err.Error())
	}
	gameAdapter := game.NewMongoAdapter(mongoCollections.Games, utcTimer)
	if err := gameAdapter.EnsureIndexes(appCtx); err != nil {
		log.Fatalf("cannot create indexes on games collection: %s", err.Error())
	}

	gameManager := game.NewManager(gameAdapter, usersAdapter, orgAdapter)

	// Echo
	valid := validator.New()
	jwtSecret := env.Env("JWT_SECRET", "jwt-token-123")
	jwtInstance := jwt.NewJWT(utcTimer, []byte(jwtSecret), time.Duration(168)*time.Hour)

	authRouter := authMux.NewMux(valid, usersAdapter, utcTimer, jwtInstance)
	orgRouter := orgMux.NewMux(valid, orgAdapter, usersAdapter)
	gameRouter := gameMux.NewMux(valid, gameManager)
	newsRouter := newsMux.NewMux(valid)

	isDebug := env.Env("DEBUG", "true")
	e := webapi.NewEcho(
		jwtInstance,
		webapi.EchoRouters{
			AuthMux:    authRouter,
			OrgRouter:  orgRouter,
			GameRouter: gameRouter,
			NewsRouter: newsRouter,
		},
		log,
		isDebug == "true")

	// Start server
	port := env.Env("APP_PORT", "8080")
	log.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
