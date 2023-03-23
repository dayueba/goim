package dao

import (
	"context"
	"time"

	"github.com/Terry-Mao/goim/internal/logic/conf"
	"github.com/gomodule/redigo/redis"
)

// Dao dao.
type Dao struct {
	c           *conf.Config
	mq          PushMsg
	redis       *redis.Pool
	redisExpire int32
}

type PushMsg interface {
	SendMessage(topic, ackInbox string, key string, msg []byte) error // ****** 这里小改了个方法名!!! 注意
	Close() error
}

// New a dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:           c,
		redis:       newRedis(c.Redis),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	d.mq = newKafka(c.Kafka)
	return d
}

func newRedis(c *conf.Redis) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: time.Duration(c.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Network, c.Addr,
				redis.DialConnectTimeout(time.Duration(c.DialTimeout)),
				redis.DialReadTimeout(time.Duration(c.ReadTimeout)),
				redis.DialWriteTimeout(time.Duration(c.WriteTimeout)),
				redis.DialPassword(c.Auth),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

// Close close the resource.
func (d *Dao) Close() error {
	return d.redis.Close()
}

// Ping dao ping.
func (d *Dao) Ping(c context.Context) error {
	return d.pingRedis(c)
}
