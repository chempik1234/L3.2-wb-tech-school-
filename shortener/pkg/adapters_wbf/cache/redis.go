package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chempik1234/super-danis-library-golang/pkg/genericports"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

// RedisWBFCache - implement genericports.GenericCachePort
type RedisWBFCache[K comparable, V genericports.ObjectWithIdentifier[K]] struct {
	client        *redis.Client
	retryStrategy retry.Strategy
}

// NewRedisWBFCache creates a new instance of RedisWBFCache
func NewRedisWBFCache[K comparable, V genericports.ObjectWithIdentifier[K]](redisClient *redis.Client, retryStrategy retry.Strategy) *RedisWBFCache[K, V] {
	return &RedisWBFCache[K, V]{client: redisClient, retryStrategy: retryStrategy}
}

// GetObjectByID - impl genericports.GenericCachePort.GetObjectByID
func (s *RedisWBFCache[K, V]) GetObjectByID(ctx context.Context, id K) (*V, error) {
	key := generateKey(id)
	data, err := s.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			return nil, nil // Not found
		}
		return nil, err // Other errors
	}

	var obj V
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

// SaveObject - impl genericports.GenericCachePort.SaveObject
func (s *RedisWBFCache[K, V]) SaveObject(ctx context.Context, fullyReadyObject *V) (*V, error) {
	key := generateKey((*fullyReadyObject).GetUniqueIdentifier())
	data, err := json.Marshal(fullyReadyObject)
	if err != nil {
		return nil, err
	}

	if err := s.client.SetWithRetry(ctx, s.retryStrategy, key, data); err != nil {
		return nil, err
	}

	return fullyReadyObject, nil
}

// DeleteObject - impl genericports.GenericCachePort.DeleteObject
func (s *RedisWBFCache[K, V]) DeleteObject(ctx context.Context, id K) error {
	key := generateKey(id)
	return s.client.Del(ctx, key)
}

// generateKey generates a Redis key based on the ID
func generateKey[K comparable](id K) string {
	return fmt.Sprintf("gnrc_rds_%v", id)
}
