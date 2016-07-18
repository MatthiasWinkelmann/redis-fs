package redisfs

import "testing"
import "github.com/hanwen/go-fuse/fuse/pathfs"
import . "github.com/smartystreets/goconvey/convey"

func TestRedisFs(t *testing.T) {
	fs := &RedisFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Host:       "127.0.0.1",
		Port:       6379,
		Db:         9,
		Dirs:       make(map[string][]string),
		Sep:        ":",
	}

	Convey("RedisFs#nameToKey()", t, func() {
		Convey("should replace '/' with ':'", func() {
			So(fs.nameToKey("foo"), ShouldEqual, "foo")
			So(fs.nameToKey("foo/bar"), ShouldEqual, "foo:bar")
			So(fs.nameToKey("foo/bar/baz"), ShouldEqual, "foo:bar:baz")
			So(fs.nameToKey("foo:bar/baz"), ShouldEqual, "foo:bar:baz")
		})
	})

	Convey("RedisFs#keyToName()", t, func() {
		Convey("should replace ':' with '/'", func() {
			So(fs.keyToName("foo"), ShouldEqual, "foo")
			So(fs.keyToName("foo:bar"), ShouldEqual, "foo/bar")
			So(fs.keyToName("foo:bar:baz"), ShouldEqual, "foo/bar/baz")
			So(fs.keyToName("foo/bar:baz"), ShouldEqual, "foo\uffffbar/baz")
		})
	})

	Convey("RedisFs#stringInSlice()", t, func() {
		Convey("check if specific string is in slice", func() {
			sliceOne := []string{"foo", "bar", "baz"}
			sliceTwo := []string{"foo", "baz"}

			resOneExist, resOneIndex := fs.stringInSlice("bar", sliceOne)
			resTwoExist, resTwoIndex := fs.stringInSlice("bar", sliceTwo)

			So(resOneExist, ShouldEqual, true)
			So(resOneIndex, ShouldEqual, 1)
			So(resTwoExist, ShouldEqual, false)
			So(resTwoIndex, ShouldEqual, -1)
		})
	})
}
