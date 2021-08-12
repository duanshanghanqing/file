package db

import (
	"file/lib/config"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool *redis.Pool
)

// 创建连接
//func createConn() (redis.Conn, error) {
//	return redis.Dial("tcp", fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort), redis.DialPassword(config.RedisPassword))
//}

// 创建链接池
func newPool(host string, port int64, password string) *redis.Pool {
	createConn := func() (redis.Conn, error) {
		return redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), redis.DialPassword(password))
	}
	return &redis.Pool{
		MaxIdle:     30,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: createConn,
	}
}

// 初始化链接池
func init() {
	pool = newPool(config.RedisHost, config.RedisPort, config.RedisPassword)
}

// GetRedisPool 获取一个链接
func GetRedisPool() redis.Conn {
	if pool == nil {
		pool = newPool(config.RedisHost, config.RedisPort, config.RedisPassword)
	}
	return pool.Get()
}
