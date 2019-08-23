package utils

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

// 生成连接池
var pool = newPool()

func CacheSet(key string, value string) bool {
	c := pool.Get()
	defer c.Close()
	c.Send("MULTI")
	c.Send("SET", key, value)
	c.Send("EXPIRE", key, 300)
	if _, err := c.Do("EXEC"); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CacheGet(key string) ([]byte, bool) {
	// 从连接池里面获得一个连接
	c := pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer c.Close()
	if cachestr, err := redis.String(c.Do("GET", key)); err == nil {
		return []byte(cachestr), true
	} else {
		log.Print(err)
	}
	return nil, false
}

// 重写生成连接池方法
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 5000, // max number of connections
		Dial: func() (redis.Conn, error) {
			//c, err := redis.Dial("tcp", "10.10.150.20:6379")
			c, err := redis.Dial("tcp", "10.10.74.170:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
