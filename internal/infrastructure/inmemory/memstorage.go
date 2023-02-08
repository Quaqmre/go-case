package inmemory

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type MemStore interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type Redis interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type storage struct {
	redisClient Redis
}

func NewStorage(redis Redis) (*storage, error) {
	s := &storage{redisClient: redis}
	return s, nil
}

func (s *storage) Set(key, value string) error {
	err := s.redisClient.Set(context.Background(), key, value, 0).Err()

	return errors.Wrapf(err, "could not set (%s: %s)", key, value)
}

func (s *storage) Get(key string) (string, error) {
	value, err := s.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "could not get the value of the key: %s", key)
	}

	return value, nil
}
