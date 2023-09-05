// Package main is an entry point to this microservice
package main

import (
	"context"
	"fmt"
	"net"

	"github.com/eugenshima/price-service/internal/config"
	"github.com/eugenshima/price-service/internal/handlers"
	"github.com/eugenshima/price-service/internal/model"
	"github.com/eugenshima/price-service/internal/repository"
	"github.com/eugenshima/price-service/internal/service"
	proto "github.com/eugenshima/price-service/proto"

	"github.com/go-playground/validator"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NewRedis function provides Connection with Redis database
func NewRedis(env string) (*redis.Client, error) {
	opt, err := redis.ParseURL(env)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis: %v", err)
	}

	fmt.Println("Connected to redis!")
	rdb := redis.NewClient(opt)
	return rdb, nil
}

// main function to run the application
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Errorf("Error extracting env variables: %v", err)
		return
	}
	client, err := NewRedis(cfg.RedisConnectionString)
	if err != nil {
		logrus.WithFields(logrus.Fields{"str": cfg.RedisConnectionString}).Errorf("NewRedis: %v", err)
	}

	pubSub := model.NewPubSub()

	r := repository.NewRedisConsumer(client)
	srv := service.NewPriceServiceService(r, pubSub)
	go srv.PublishToAllSubscribers(context.Background())
	hndl := handlers.NewPriceServiceHandler(srv, validator.New())

	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		logrus.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	proto.RegisterPriceServiceServer(serverRegistrar, hndl)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		logrus.Fatalf("cannot start server: %s", err)
	}
}
