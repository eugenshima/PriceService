// Package main is an ebtry point to this microservice
package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/eugenshima/PriceService/internal/config"
	"github.com/eugenshima/PriceService/internal/handlers"
	"github.com/eugenshima/PriceService/internal/repository"
	"github.com/eugenshima/PriceService/internal/service"
	proto "github.com/eugenshima/PriceService/proto"

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
	var mu sync.RWMutex
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Errorf("Error extracting env variables: %v", err)
		return
	}
	client, err := NewRedis(cfg.RedisConnectionString)
	if err != nil {
		logrus.WithFields(logrus.Fields{"str": cfg.RedisConnectionString}).Errorf("NewRedis: %v", err)
	}

	r := repository.NewRedisConsumer(client)
	srv := service.NewPriceServiceService(r)
	hndl := handlers.NewPriceServiceHandler(srv, &mu)

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
