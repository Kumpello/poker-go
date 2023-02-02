package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"pokergo/internal/articles"
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
	orgAdapter := org.NewMongoAdapter(mongoCollections.Org, utcTimer)
	gameAdapter := game.NewMongoAdapter(mongoCollections.Games, utcTimer)
	artsAdapter := articles.NewMongoAdapter(mongoCollections.Arts)
	gameManager := game.NewManager(gameAdapter, usersAdapter, orgAdapter)

	// Echo
	jwtSecret := env.Env("JWT_SECRET", "jwt-token-123")
	jwtInstance := jwt.NewJWT(utcTimer, []byte(jwtSecret), time.Duration(168)*time.Hour)
	authRouter := authMux.NewMux(usersAdapter, utcTimer, jwtInstance)
	orgRouter := orgMux.NewMux(orgAdapter, usersAdapter)
	gameRouter := gameMux.NewMux(gameManager)
	newsRouter := newsMux.NewMux(artsAdapter)

	isDebug := env.Env("DEBUG", "true")
	validate := validator.New()
	e := webapi.NewEcho(
		validate,
		jwtInstance,
		webapi.EchoRouters{
			AuthRouter: authRouter,
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
