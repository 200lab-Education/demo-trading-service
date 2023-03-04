package redlock

import "context"

func (r *redisLock) UnLock(ctx context.Context, key string) error {
	if key == "" {
		return KeyNullErr
	}

	for _, client := range r.clients {
		exsited := client.Get(ctx, key)
		if exsited == nil {
			return KeyNotExistedErr
		}

		_, err := client.Del(ctx, key).Result()
		if err != nil {
			return err
		}
	}

	return nil
}
