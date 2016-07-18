package redisfs

import "os"
import "log"
import "path"
import "regexp"
import "strings"
import "github.com/hanwen/go-fuse/fuse"
import "github.com/garyburd/redigo/redis"
import "github.com/hanwen/go-fuse/fuse/nodefs"
import "github.com/hanwen/go-fuse/fuse/pathfs"

type RedisFs struct {
	pathfs.FileSystem
	Host string
	Port int
	Auth string
	Dirs map[string][]string
	Db   int
	Sep  string
	pool *redis.Pool
}

func (fs *RedisFs) Init() {
	pool := &redis.Pool{
		MaxIdle:   2,
		MaxActive: 20,
		Dial: func() (redis.Conn, error) {
			return fs.CreateRedisConn()
		},
	}

	fs.pool = pool
}

func (fs *RedisFs) CreateRedisConn() (redis.Conn, error) {
	return NewRedisConn(fs.Host, fs.Port, fs.Db, fs.Auth)
}

func (fs *RedisFs) GetAttr(name string, ctx *fuse.Context) (*fuse.Attr, fuse.Status) {
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	// ignore hidden files
	if string(name[0]) == "." {
		return nil, fuse.ENOENT
	}

	// find dir in memory
	dirs, ok := fs.Dirs[path.Dir(name)]
	baseName := path.Base(name)

	if ok {
		exist, _ := fs.stringInSlice(baseName, dirs)
		if exist {
			return &fuse.Attr{
				Mode: fuse.S_IFDIR | 0755,
			}, fuse.OK
		}
	}

	// Open connection
	conn := fs.pool.Get()
	defer conn.Close()

	// find attr in redis
	key := fs.nameToKey(name)
	content, err1 := redis.String(conn.Do("GET", key))
	list, err2 := redis.Strings(conn.Do("KEYS", key+fs.Sep+"*"))

	switch {
	case err2 == nil && len(list) > 0:
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
		break
	case err1 == nil:
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644,
			Size: uint64(len(content)),
		}, fuse.OK
		break
	}

	return nil, fuse.ENOENT
}

func (fs *RedisFs) OpenDir(name string, ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	conn := fs.pool.Get()
	defer conn.Close()

	pattern := fs.nameToPattern(name)
	res, err := redis.Strings(conn.Do("KEYS", pattern))

	if err != nil {
		fs.printError(err)
		return nil, fuse.ENOENT
	}

	m := make(map[string]bool)
	entries := append(
		fs.dirsToEntries(name, m),
		fs.resToEntries(fs.nameToKey(name), res, m)...)

	return entries, fuse.OK
}

func (fs *RedisFs) Open(name string, flags uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	conn := fs.pool.Get()
	defer conn.Close()

	key := fs.nameToKey(name)
	_, err := conn.Do("EXISTS", key)

	if err != nil {
		fs.printError(err)
		return nil, fuse.ENOENT
	}

	return NewRedisFile(fs.pool, key), fuse.OK
}

func (fs *RedisFs) Create(name string, flags uint32, mode uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	conn := fs.pool.Get()
	defer conn.Close()

	key := fs.nameToKey(name)
	_, err := conn.Do("SET", key, "")

	if err != nil {
		fs.printError(err)
		return nil, fuse.ENOENT
	}

	return NewRedisFile(fs.pool, key), fuse.OK
}

func (fs *RedisFs) Rename(oldName string, newName string, ctx *fuse.Context) fuse.Status {
	oldKey := fs.nameToKey(oldName)
	newKey := fs.nameToKey(newName)

	conn := fs.pool.Get()
	defer conn.Close()

	// get file content
	content, err := redis.String(conn.Do("GET", oldKey))

	if err != nil {
		return fuse.ENOENT
	}

	// create new file
	_, err = conn.Do("SET", newKey, content)

	if err != nil {
		return fuse.ENOENT
	}

	// delete old file
	_, err = conn.Do("DEL", oldKey)

	if err != nil {
		return fuse.ENOENT
	}

	return fuse.OK
}

func (fs *RedisFs) Unlink(name string, ctx *fuse.Context) fuse.Status {
	if name == "" {
		return fuse.OK
	}

	conn := fs.pool.Get()
	defer conn.Close()

	key := fs.nameToKey(name)
	_, err := conn.Do("DEL", key)

	if err != nil {
		fs.printError(err)
		return fuse.ENOENT
	}

	return fuse.OK
}

func (fs *RedisFs) Rmdir(name string, ctx *fuse.Context) fuse.Status {
	if name == "" {
		return fuse.OK
	}

	// check if name is in memory
	dirName := path.Dir(name)
	dir, ok := fs.Dirs[dirName]
	baseName := path.Base(name)

	if ok {
		exist, index := fs.stringInSlice(baseName, dir)
		if exist {
			fs.Dirs[dirName] = append(dir[:index], dir[index+1:]...)
			return fuse.OK
		}
	}

	// open connection
	conn := fs.pool.Get()
	defer conn.Close()

	// if name isn't in memory then find it in redis
	pattern := fs.nameToPattern(name)
	list, err := redis.Strings(conn.Do("KEYS", pattern))

	if err != nil {
		fs.printError(err)
		return fuse.ENOENT
	}

	for _, el := range list {
		_, err := conn.Do("DEL", el)
		if err != nil {
			fs.printError(err)
			return fuse.ENOENT
		}
	}

	return fuse.OK
}

func (fs *RedisFs) Mkdir(name string, mode uint32, ctx *fuse.Context) fuse.Status {
	dir := path.Join(name, "..")

	_, ok := fs.Dirs[dir]

	if !ok {
		fs.Dirs[dir] = make([]string, 0, 10)
	}

	fs.Dirs[dir] = append(fs.Dirs[dir], path.Base(name))

	return fuse.OK
}

func (fs *RedisFs) nameToPattern(name string) string {
	pattern := fs.nameToKey(name)

	if name == "" {
		pattern += "*"
	} else {
		pattern += fs.Sep + "*"
	}

	return pattern
}

func (fs *RedisFs) dirsToEntries(dir string, m map[string]bool) []fuse.DirEntry {
	entries := make([]fuse.DirEntry, 0, 2)

	if dir == "" {
		dir = "."
	}

	if list, ok := fs.Dirs[dir]; ok {
		for _, key := range list {
			m[key] = true
			entries = append(entries, fuse.DirEntry{
				Name: key,
				Mode: fuse.S_IFDIR,
			})
		}
	}

	return entries
}

func (fs *RedisFs) resToEntries(dir string, list []string, m map[string]bool) []fuse.DirEntry {
	entries := make([]fuse.DirEntry, 0, 2)
	offset := len(dir)
	sepCount := strings.Count(dir, string(os.PathSeparator)) + 1

	if offset != 0 {
		offset += 1
	}

	for _, el := range list {
		key := el[offset:]
		keySepCount := strings.Count(key, fs.Sep)

		switch true {
		case keySepCount == 0:
			entries = append(entries, fuse.DirEntry{
				Name: fs.keyToName(key),
				Mode: fuse.S_IFREG,
			})
			break
		case keySepCount >= sepCount:
			tmp := strings.SplitN(key, fs.Sep, 2)
			key = tmp[0]

			if _, ok := m[key]; !ok {
				m[key] = true
				entries = append(entries, fuse.DirEntry{
					Name: key,
					Mode: fuse.S_IFDIR,
				})
			}
		}
	}

	return entries
}

func (fs *RedisFs) nameToKey(name string) string {
	re := regexp.MustCompile(string(os.PathSeparator))
	key := re.ReplaceAllLiteralString(name, fs.Sep)
	key = fs.decodePathSeparator(key)
	return key
}

func (fs *RedisFs) keyToName(key string) string {
	name := fs.encodePathSeparator(key)
	re := regexp.MustCompile(fs.Sep)
	name = re.ReplaceAllLiteralString(name, string(os.PathSeparator))
	return name
}

func (fs *RedisFs) encodePathSeparator(str string) string {
	re := regexp.MustCompile(string(os.PathSeparator))
	str = re.ReplaceAllLiteralString(str, "\uffff")
	return str
}

func (fs *RedisFs) decodePathSeparator(str string) string {
	re := regexp.MustCompile("\uffff")
	str = re.ReplaceAllLiteralString(str, string(os.PathSeparator))
	return str
}

func (fs *RedisFs) printError(err error) {
	log.Println("Error:", err)
}

func (fs *RedisFs) stringInSlice(target string, list []string) (bool, int) {
	for i, str := range list {
		if str == target {
			return true, i
		}
	}
	return false, -1
}
