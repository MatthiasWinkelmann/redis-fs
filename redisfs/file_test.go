package redisfs

import "time"
import "testing"
import "github.com/hanwen/go-fuse/fuse"
import "github.com/garyburd/redigo/redis"
import . "github.com/smartystreets/goconvey/convey"

func TestRedisFile(t *testing.T) {
	pool := &redis.Pool{
		Dial: dialTestDB,
	}

	defer func() {
		conn := pool.Get()
		conn.Do("FLUSHALL")
		conn.Close()
		pool.Close()
	}()

	Convey("Write", t, func() {
		Convey("should work", func() {
			conn := pool.Get()
			defer conn.Close()

			data := []byte("Ghost Island Taiwan")

			_, err := conn.Do("SET", "writing", "")

			file := NewRedisFile(pool, "writing")
			_, code := file.Write(data, 0)

			So(code, ShouldEqual, fuse.OK)
			res, err := redis.String(conn.Do("GET", "writing"))

			if err != nil {
				panic(err)
			}

			So(res, ShouldEqual, string(data))
		})
	})

	Convey("Read", t, func() {
		Convey("should work", func() {
			conn := pool.Get()
			defer conn.Close()

			file := NewRedisFile(pool, "reading")
			data := []byte("QQ")
			_, err := conn.Do("SET", "reading", string(data))

			if err != nil {
				panic(err)
			}

			buf := make([]byte, 100)
			res, code := file.Read(buf, 0)

			So(code, ShouldEqual, fuse.OK)
			So(res.Size(), ShouldEqual, len(data))
		})
	})
}

func dialTestDB() (redis.Conn, error) {
	c, err := redis.DialTimeout("tcp", ":6379", 0, 1*time.Second, 1*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = c.Do("SELECT", "9")
	if err != nil {
		return nil, err
	}

	return c, nil
}
