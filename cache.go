package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	redisURL := os.Getenv("REDIS_URL")
	pool = newPool(redisURL)
}

func GetCacheItem(identifier string) *Item {
	conn := pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("GET", identifier))
	if err == redis.ErrNil {
		return nil
	}
	if err != nil {
		log.Fatalf("Failed to get item from cache - %s: %s", identifier, err)
		return nil
	}
	jsonBody := []byte(res)

	var item Item
	err = json.Unmarshal(jsonBody, &item)
	if err != nil {
		log.Fatalf("Failed to decode the value as json - %s: %s", res, err)
		return nil
	}
	return &item
}

func SetCacheItem(identifier string, item *Item) {
	if item == nil {
		return
	}

	b, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("Failed to encode the value as json - %s: %s", identifier, err)
		return
	}

	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", identifier, 300, string(b))
	if err != nil {
		log.Fatalf("Failed to set item in redis - %s: %s", identifier, err)
		return
	}
}

func newPool(redisURL string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisURL)
		},
	}
}
