package dao

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

const password = "root"
const host = "127.0.0.1:6379"
const redisDB = 0

var Pool *redis.Pool
func InitRedis() {

	Pool = &redis.Pool{
		MaxIdle: 20,
		IdleTimeout: 60 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_,err := c.Do("ping")

			return err
		},
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			_, err = c.Do("auth", password)
			if err != nil {
				return nil, err
			}

			_, err = c.Do("select", redisDB)
			return c, err
		},
	}
}
