package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Quaqmre/go-case/internal/infrastructure/inmemory"
	"github.com/Quaqmre/go-case/internal/infrastructure/persistent"
	"github.com/Quaqmre/go-case/internal/router"
	"github.com/Quaqmre/go-case/internal/util"
	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var logger log.Logger

func init() {

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	lvl := util.GetLogLevelFromEnv()
	logger = level.NewFilter(logger, lvl)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		level.Debug(logger).Log("err", "Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		level.Warn(logger).Log("err", fmt.Sprintf("error in reading port\n: %s", err.Error()))
	}
	appAddress := fmt.Sprintf("0.0.0.0:%d", port)

	//Storage
	redisClient, err := ConnectRedis()
	ctx, mongoClient, collection, err := ConnectMongoDb()

	if err != nil {
		level.Warn(logger).Log("err", fmt.Sprintf("error in connect redis\n: %s", err.Error()))
	}

	//Service init
	inmemStorage, err := inmemory.NewStorage(redisClient)
	if err != nil {
		level.Warn(logger).Log("err", fmt.Sprintf("error in create redis storage\n: %s", err.Error()))
	}
	mongoStorage, err := persistent.NewDb(ctx, mongoClient, collection)

	if err != nil {
		level.Warn(logger).Log("err", fmt.Sprintf("error in create mongo storage\n: %s", err.Error()))
	}

	r := router.RegisterRoutes(inmemStorage, mongoStorage, logger)

	srv := &http.Server{
		Addr:    appAddress,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			level.Warn(logger).Log(err)
		}
	}()
	level.Info(logger).Log("err", fmt.Sprintf("server started... %s", appAddress))

	shutdownGracefully(srv)

}

func shutdownGracefully(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	level.Info(logger).Log("err", "server stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		level.Warn(logger).Log(err)
	}
}
func ConnectRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, errors.Wrap(err, "invalid Redis url")
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to Redis")
	}

	return client, nil
}

func ConnectMongoDb() (context.Context, *mongo.Client, *mongo.Collection, error) {
	uri := os.Getenv("MONGO_URI")

	ctx := context.TODO()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "database connection error")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to ping the primary")
	}

	dbName := os.Getenv("MONGO_DB")
	database := client.Database(dbName)

	collectionName := os.Getenv("MONGO_COLLECTION")
	collection := database.Collection(collectionName)

	return ctx, client, collection, nil
}
