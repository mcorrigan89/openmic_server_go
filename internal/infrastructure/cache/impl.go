package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func CreateCache() *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Cache{
		client: rdb,
	}
}

func (cache *Cache) Get(ctx context.Context, key string) (*string, error) {
	val, err := cache.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &val, nil
	}
}

func (cache *Cache) Set(ctx context.Context, key, value string, exp time.Duration) error {
	err := cache.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) Close() {
	cache.client.Close()
}

func (cache *Cache) Del(ctx context.Context, key string) error {
	err := cache.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
