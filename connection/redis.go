package connection

import "github.com/go-redis/redis/v8"

func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})
}
