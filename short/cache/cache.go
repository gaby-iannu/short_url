package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)


type Cache interface {
	Put(string, string) error 
	Get(string) (string, error)
}

type cache struct {
	client *redis.Client
}

type error interface {
	Error() string
}

type NotExistError struct {

}

func (e *NotExistError) Error() string {
	return fmt.Sprintf("Value dosn't exit into cache")
}

func New(url string, port int) Cache {
	return &cache{
		client: redis.NewClient(&redis.Options{
			Addr:fmt.Sprintf("%s:%d", url, port),
		}),
	}
}

func (c *cache) Put(key, value string) error {
	 return c.client.Set(key, value, 0).Err()
}

func (c *cache) Get(key string) (string, error) {
	cmd := c.client.Get(key)

	value, err := cmd.Result()
	
	if err == redis.Nil {
		return "",&NotExistError{}
	}

	return value, nil
}
