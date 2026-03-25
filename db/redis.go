package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	_redisConn *redis.Client
	once       sync.Once
)

func RedisConnect() (*redis.Client, error) {
	once.Do(func() {
		addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
		poolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOLSIZE"))

		log.Printf("Connecting to Redis %s ...\n", addr)

		_redisConn = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
			PoolSize: poolSize,
		})

		_, err := _redisConn.Ping(context.Background()).Result()
		if err != nil {
			log.Printf("Redis connection failed: %v\n", err)
			_redisConn = nil
		}
	})

	return _redisConn, nil
}

func RedisGetConnection() *redis.Client {
	conn, _ := RedisConnect()
	return conn
}

func RedisClose() {
    if _redisConn != nil {
        if err := _redisConn.Close(); err != nil {
            log.Printf("Error closing Redis connection: %v\n", err)
        }
    }
}