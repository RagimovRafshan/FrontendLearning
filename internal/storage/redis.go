package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"messenger/config"
)

// RedisClient обертка над redis.Client для работы с Redis
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient создает новый клиент Redis
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Проверка подключения
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// Client возвращает базовый клиент redis
func (c *RedisClient) Client() *redis.Client {
	return c.client
}

// Close закрывает соединение с Redis
func (c *RedisClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// HealthCheck проверяет подключение к Redis
func (c *RedisClient) HealthCheck(ctx context.Context) error {
	result, err := c.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	if result != "PONG" {
		return fmt.Errorf("unexpected Redis ping response: %s", result)
	}
	return nil
}

// Set устанавливает значение ключа с опциональным временем жизни
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get получает значение по ключу
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete удаляет ключ
func (c *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists проверяет существование ключа
func (c *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire устанавливает время жизни для ключа
func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL получает оставшееся время жизни ключа
func (c *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Pub/Sub методы

// Publish публикует сообщение в канал
func (c *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe подписывается на канал
func (c *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

// List методы

// LPush добавляет элемент в начало списка
func (c *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.LPush(ctx, key, values...).Err()
}

// RPush добавляет элемент в конец списка
func (c *RedisClient) RPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.RPush(ctx, key, values...).Err()
}

// LRange получает диапазон элементов из списка
func (c *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LLen получает длину списка
func (c *RedisClient) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

// Hash методы

// HSet устанавливает значение поля в хеше
func (c *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.client.HSet(ctx, key, values...).Err()
}

// HGet получает значение поля из хеша
func (c *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll получает все поля и значения хеша
func (c *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel удаляет поле из хеша
func (c *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// Sorted Set методы

// ZAdd добавляет элемент в сортированное множество
func (c *RedisClient) ZAdd(ctx context.Context, key string, score float64, member interface{}) error {
	return c.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// ZRange получает диапазон элементов из сортированного множества
func (c *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRem удаляет элемент из сортированного множества
func (c *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

// ZCard получает количество элементов в сортированном множестве
func (c *RedisClient) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, key).Result()
}
