package redis

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	er "github.com/swethabhageerath/redis/lib/errors"
	"github.com/swethabhageerath/redis/lib/models"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) Redis {
	return Redis{
		client: client,
	}
}

func (r Redis) Set(key string, value string, expiration time.Duration, out chan error) {
	s := r.client.Set(key, value, expiration)
	err := s.Err()
	if err != nil {
		out <- errors.Wrap(err, er.ErrSetRedisCache.String())
	}
}

func (r Redis) Get(key string, out chan models.RedisGetResponse) {
	s := r.client.Get(key)
	result, err := s.Result()
	if err != nil {
		if err == redis.Nil {
			out <- models.RedisGetResponse{
				Data:  "",
				Error: errors.Wrap(err, fmt.Sprintf(er.ErrKeyNotExists.String(), key)),
			}
		} else {
			out <- models.RedisGetResponse{
				Data:  "",
				Error: errors.Wrap(err, fmt.Sprintf(er.ErrRetrievingKey.String(), key)),
			}
		}
	}
	out <- models.RedisGetResponse{
		Data:  result,
		Error: nil,
	}
}

func (r Redis) Remove(key string, expiration time.Duration, out chan error) {
	e := r.client.Expire(key, expiration)
	err := e.Err()
	if err != nil {
		if err == redis.Nil {
			out <- errors.Wrap(err, fmt.Sprintf(er.ErrKeyNotExists.String(), key))
		} else {
			out <- errors.Wrap(err, fmt.Sprintf(er.ErrRetrievingKey.String(), key))
		}
	}

	out <- nil
}
