package redis

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/swethabhageerath/logging/lib/constants"
	m "github.com/swethabhageerath/logging/lib/models"
	w "github.com/swethabhageerath/logging/lib/writers"
	E "github.com/swethabhageerath/redis/lib/errors"
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
		er := errors.Wrap(err, E.ErrSetRedisCache.String())
		r.log(er)
		out <- er
	}
}

func (r Redis) Get(key string, out chan models.RedisGetResponse) {
	s := r.client.Get(key)
	result, err := s.Result()
	if err != nil {
		if err == redis.Nil {
			er := errors.Wrap(err, fmt.Sprintf(E.ErrKeyNotExists.String(), key))
			r.log(er)
			out <- models.RedisGetResponse{
				Data:  "",
				Error: er,
			}
		} else {
			er := errors.Wrap(err, fmt.Sprintf(E.ErrRetrievingKey.String(), key))
			r.log(er)
			out <- models.RedisGetResponse{
				Data:  "",
				Error: er,
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
		var er error
		if err == redis.Nil {
			er = errors.Wrap(err, fmt.Sprintf(E.ErrKeyNotExists.String(), key))
			r.log(er)
			out <- er
		} else {
			er = errors.Wrap(err, fmt.Sprintf(E.ErrRetrievingKey.String(), key))
			r.log(er)
			out <- er
		}
	}

	out <- nil
}

func (r Redis) log(err error) {
	l := m.New(m.WithMandatoryFields("Redis", "bmoola", constants.ERROR), m.WithRequestId("abc123"), m.WithStackTrace(err))
	_, e := l.Attach(w.FileWriter{})
	if e != nil {
		panic(e)
	}
	l.Notify()
}
