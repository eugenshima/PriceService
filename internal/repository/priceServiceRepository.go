// Package repository provides functions for interacting with a redis stream
package repository

import (
	"context"
	"encoding/json"

	"github.com/eugenshima/PriceService/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisConsumer is a struct for Redis Stream Consumer
type RedisConsumer struct {
	redisClient *redis.Client
}

// NewConsumer creates a new Redis Stream Consumer
func NewRedisConsumer(redisClient *redis.Client) *RedisConsumer {
	return &RedisConsumer{redisClient: redisClient}
}

// RedisConsumer gets latest price from redis stream
func (repo *RedisConsumer) RedisConsumer(ctx context.Context) ([]*model.Share, error) {
	var shares []*model.Share

	streams := repo.redisClient.XRevRange(ctx, "PriceStreaming", "+", "-").Val()

	err := json.Unmarshal([]byte(streams[0].Values["GeneratedPrices"].(string)), &shares)
	if err != nil {
		logrus.WithFields(logrus.Fields{"shares": shares}).Errorf("Error unmarshalling: %v", err)
		return nil, err
	}
	return shares, nil
}
