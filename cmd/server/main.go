package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"pokergo/internal/mongo"
	"pokergo/internal/org"
	"pokergo/internal/users"
	"pokergo/internal/webapi"
	"pokergo/internal/webapi/auth"
	webOrg "pokergo/internal/webapi/org"
	"pokergo/pkg/env"
	"pokergo/pkg/jwt"
	"pokergo/pkg/logger"
	"pokergo/pkg/timer"
	"time"
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
	userAdapter := users.NewMongoAdapter(mongoCollections.Users, log)
	if err := userAdapter.EnsureIndexes(appCtx); err != nil {
		log.Fatalf("cannot create indexes on users collection", err)
	}
	orgAdapter := org.NewMongoAdapter(mongoCollections.Org, utcTimer)
	if err := userAdapter.EnsureIndexes(appCtx); err != nil {
		log.Fatalf("cannot create indexes on organizations collection", err)
	}

	// Echo
	valid := validator.New()
	jwtSecret := env.Env("JWT_SECRET", "jwt-token-123")
	jwtInstance := jwt.NewJWT(utcTimer, []byte(jwtSecret), time.Duration(168)*time.Hour)

	authMux := auth.NewMux(userAdapter, utcTimer, jwtInstance, valid)
	orgMux := webOrg.NewMux(orgAdapter, userAdapter)

	e := webapi.NewEcho(jwtInstance, authMux, orgMux)

	// Start server
	port := env.Env("APP_PORT", "8080")
	log.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
