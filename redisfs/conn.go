package redisfs

import "log"
import "strconv"
import "github.com/garyburd/redigo/redis"

func NewRedisConn(host string, port int, db int, auth string) (redis.Conn, error) {
	address := host + ":" + strconv.Itoa(port)
	conn, err := redis.Dial("tcp", address)

	if err != nil {
		return nil, err
	}

	if len(auth) > 0 {
		log.Println("Info:", "Auth")
		if _, err := conn.Do("AUTH", auth); err != nil {
			conn.Close()
			return nil, err
		}
	}

	if db != 0 {
		log.Println("Info:", "Use DB", db)
		if _, err := conn.Do("SELECT", db); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
