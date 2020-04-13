package db

import "github.com/go-redis/redis"

var redisClient *redis.Client

// GetRedisClient returns a redis connection client
func GetRedisClient() (*redis.Client, error) {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		// connection test
		_, err := redisClient.Ping().Result()
		if err != nil {
			return nil, err
		}
	}
	return redisClient, nil
}
