package redisfs

import (
	"bytes"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type redisFile struct {
	pool *redis.Pool
	key  string
}

func NewRedisFile(pool *redis.Pool, key string) nodefs.File {
	file := new(redisFile)
	file.pool = pool
	file.key = key
	return file
}

func (f *redisFile) SetInode(*nodefs.Inode) {
}

func (f *redisFile) InnerFile() nodefs.File {
	return nil
}

func (f *redisFile) String() string {
	return "redisFile"
}

func (f *redisFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	conn := f.pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", f.key))

	if err != nil {
		log.Println("ERROR:", err)
		return nil, fuse.EIO
	}

	end := int(off) + int(len(buf))
	dataLen := len(data)

	if end > dataLen {
		end = dataLen
	}

	return fuse.ReadResultData(data[off:end]), fuse.OK
}

func (f *redisFile) Flock(int) fuse.Status {
	return fuse.OK
}
func (f *redisFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	conn := f.pool.Get()
	defer conn.Close()

	originalData, err := redis.Bytes(conn.Do("GET", f.key))

	if err != nil {
		log.Println("Error:", err)
		return 0, fuse.EIO
	}

	leftChunk := originalData[:off]
	end := int(off) + int(len(data))

	var rightChunk []byte

	if end > len(originalData) {
		rightChunk = []byte{}
	} else {
		rightChunk = data[int(off)+len(data):]
	}

	newValue := bytes.NewBuffer(leftChunk)
	newValue.Grow(len(data) + len(rightChunk))
	newValue.Write(data)
	newValue.Write(rightChunk)

	_, err = conn.Do("SET", f.key, newValue.String())

	if err != nil {
		log.Println("Error:", err)
		return 0, fuse.EIO
	}

	return uint32(len(data)), fuse.OK
}

func (f *redisFile) Flush() fuse.Status {
	return fuse.OK
}

func (f *redisFile) Release() {
}

func (f *redisFile) GetAttr(out *fuse.Attr) fuse.Status {
	conn := f.pool.Get()
	defer conn.Close()

	content, err := redis.String(conn.Do("GET", f.key))

	if err != nil {
		log.Println("Error:", err)
		return fuse.EIO
	}

	out.Mode = fuse.S_IFREG | 0644
	out.Size = uint64(len(content))

	return fuse.OK
}

func (f *redisFile) GetLk(owner uint64, lk *fuse.FileLock, flags uint32, out *fuse.FileLock) (code fuse.Status) {
	return fuse.ENOSYS
}

func (f *redisFile) SetLk(owner uint64, lk *fuse.FileLock, flags uint32) (code fuse.Status) {
	return fuse.ENOSYS
}

func (f *redisFile) SetLkw(owner uint64, lk *fuse.FileLock, flags uint32) (code fuse.Status) {
	return fuse.ENOSYS
}

func (f *redisFile) Fsync(flags int) (code fuse.Status) {
	return fuse.OK
}

func (f *redisFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	return fuse.ENOSYS
}

func (f *redisFile) Truncate(size uint64) fuse.Status {
	return fuse.OK
}

func (f *redisFile) Chown(uid uint32, gid uint32) fuse.Status {
	return fuse.ENOSYS
}

func (f *redisFile) Chmod(perms uint32) fuse.Status {
	return fuse.ENOSYS
}

func (f *redisFile) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	return fuse.OK
}
