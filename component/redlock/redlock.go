package redlock

import "github.com/go-redis/redis/v8"

type redisLock struct {
	clients []*redis.Client
}

func NewRedisLock(clients []*redis.Client) *redisLock {
	return &redisLock{
		clients: clients,
	}
}
