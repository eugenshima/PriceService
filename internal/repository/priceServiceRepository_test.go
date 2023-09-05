package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/eugenshima/price-service/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var redisConn *RedisConsumer

var shares = []*model.Share{
	{
		ShareName:  "test",
		SharePrice: "test",
	},
}

func TestRedisConsumer(t *testing.T) {
	id, err := addTestingRow()
	require.NoError(t, err)
	result, err := redisConn.RedisConsumer(context.Background())
	t.Logf("Got: %v", result)
	require.NotNil(t, &result)
	require.NoError(t, err)
	err = delTestingRow("PriceStreaming", id)
	require.NoError(t, err)
}

func TestNilRedisConsumer(t *testing.T) {
	result, err := redisConn.RedisConsumer(context.Background())
	require.NotNil(t, &result)
	require.NoError(t, err)
}

func addTestingRow() (string, error) {
	jsonValue, _ := json.Marshal(shares)
	id := strconv.FormatInt(time.Now().Unix(), 10)
	err := redisConn.redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "PriceStreaming",
		ID:     id,
		Values: map[string]interface{}{
			"GeneratedPrices": jsonValue,
		},
	}).Err()
	return id, err
}

func delTestingRow(stream, id string) error {
	return redisConn.redisClient.XDel(context.Background(), stream, id).Err()
}
