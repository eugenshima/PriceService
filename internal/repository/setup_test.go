package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/redis/go-redis/v9"
)

func SetupTestRedis() (*redis.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}

	resource, err := pool.Run("redis", "latest", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	client, err := redis.ParseURL("redis://:@localhost:6379/1")
	if err != nil {
		return nil, nil, fmt.Errorf("Could not parse redis url: %s", err)
	}
	rdb := redis.NewClient(client)
	cleanup := func() {
		rdb.Close()
		pool.Purge(resource)
	}
	return rdb, cleanup, nil
}

func TestMain(m *testing.M) {
	rdb, cleanupRedis, err := SetupTestRedis()
	if err != nil {
		fmt.Println(err)
		cleanupRedis()
		os.Exit(1)
	}
	redisConn = NewRedisConsumer(rdb)
	exitVal := m.Run()
	cleanupRedis()
	os.Exit(exitVal)
}
