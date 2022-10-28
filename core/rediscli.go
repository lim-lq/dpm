package core

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/core/log"
)

type redisClient struct {
	ctx context.Context
	cli *redis.Client
}

var rediscli *redisClient

func InitRedis() {
	rHost := config.GetString("redis.host")
	if rHost == "" {
		log.Logger.Fatal("Please configure redis.host")
	}
	rPort := config.GetInt("redis.port")
	if rPort == 0 {
		log.Logger.Fatal("Please configure redis.port")
	}
	rDb := config.GetInt("redis.db")
	rPass := config.GetString("redis.pass")
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rHost, rPort),
		Password: rPass,
		DB:       rDb,
	})
	rediscli = &redisClient{
		ctx: context.Background(),
		cli: db,
	}
	if err := rediscli.Ping(); err != nil {
		log.Logger.Fatalf("Connect redis error: %v", err)
	}
}

func GetRedisClient() *redisClient {
	return rediscli
}

func (r *redisClient) Ping() error {
	_, err := r.cli.Ping(r.ctx).Result()
	return err
}

func (r *redisClient) Get(key string) (value string, err error) {
	value, err = r.cli.Get(r.ctx, key).Result()
	return
}

func (r *redisClient) Set(key string, value interface{}, timeout time.Duration) error {
	_, err := r.cli.Set(r.ctx, key, value, timeout).Result()
	return err
}

func (r *redisClient) Del(key string) error {
	_, err := r.cli.Del(r.ctx, key).Result()
	return err
}
