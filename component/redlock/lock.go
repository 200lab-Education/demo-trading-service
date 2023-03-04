package redlock

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var KeyNullErr = errors.New("KeyNull")
var KeyExistedErr = errors.New("KeyExistedErr")
var KeyNotExistedErr = errors.New("KeyNotExistedErr")
var CanNotSetKeyErr = errors.New("CanNotSetKeyErr")

func (r *redisLock) Lock(ctx context.Context, key string, value string, expire time.Duration) error {
	if key == "" || value == "" {
		return KeyNullErr
	}

	for _, client := range r.clients {
		var status string
		if err := client.Get(ctx, key).Scan(&status); err != nil && err != redis.Nil {
			return err
		}

		if status == value {
			return KeyExistedErr
		}

		if err := client.SetNX(ctx, key, value, expire).Err(); err != nil {
			return err
		}
	}

	return nil
}
