package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = errors.New("key doesn't")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val any, exp time.Duration) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{
		client: client,
	}
}

func (s *cache) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("redis get key %s error: %w", key, err)
	}

	if err != nil && errors.Is(err, redis.Nil) {
		return "", ErrKeyNotExist
	}

	return val, nil
}

func (s *cache) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	if err := s.client.Set(ctx, key, val, exp).Err(); err != nil {
		return fmt.Errorf("redis set key %s error: %w", key, err)
	}

	return nil
}
