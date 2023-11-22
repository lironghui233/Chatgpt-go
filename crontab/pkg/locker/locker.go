package locker

import (
	"context"
	"crontab/pkg/log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker interface {
	Lock(key string)
	Unlock(key string)
}

// 利用redis构建分布式锁
type redisLocker struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisLocker(client *redis.Client, ttl time.Duration) Locker {
	return &redisLocker{
		client: client,
		ttl:    ttl,
	}
}

func (l *redisLocker) Lock(key string) {
	for {
		ok, err := l.client.SetNX(context.Background(), key, "", l.ttl).Result()
		if err != nil || !ok {
			time.Sleep(500 * time.Second)
			continue
		}
		break
	}
}
func (l *redisLocker) Unlock(key string) {
	err := l.client.Del(context.Background(), key).Err()
	if err != nil {
		log.Error(err)
	}
}
