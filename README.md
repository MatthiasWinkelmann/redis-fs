redis-fs
===========

redis-fs lets you mount a Redis database as a filesystem. It is a continuation of redis-mount by Po-Ying Chen which was deleted from Github for unkown reasons.

[![Build Status](https://travis-ci.org/MatthiasWinkelmann/redis-fs.svg?branch=master)](https://travis-ci.org/MatthiasWinkelmann/redis-fs)

## Help

Please star this repository if you='re using redis-fs so that I know it's worth to contiue working on it.

Please report any Issue you may encounter using redis-fs.

Pull requests are welcome.

## Usage

```bash
redis-fs 0.2.0
$ redis-fs ~/redis

--host, -h   localhost    Redis host name
--port, -p   6379         Redis port number
--auth, -a                Redis password
--sep, -s    :            Redis key separator
```

## What we can do with it?

1. Create a fast, auto-expanding RAM disk for (i. e.) working with temporary files.
1. Use `grep` to search for text in redis values.
2. Pass data to other programs. ex: `$ cat redis-key | pretty-print`

![screenshot](documentation/screenshot.gif)

## Installation

### Download binary file

* [mac-amd64](https://github.com/MatthiasWinkelmann/redis-fs/releases/download/0.2.0/redis-fs-darwin-amd64)
* [linux-amd64](https://github.com/MatthiasWinkelmann/redis-fs/releases/download/0.2.0/redis-fs-linux-amd64)
* [linux-386](https://github.com/MatthiasWinkelmann/redis-fs/releases/download/0.2.0/redis-fs-linux-386)
* [linux-arm](https://github.com/MatthiasWinkelmann/redis-fs/releases/download/0.2.0/redis-fs-linux-arm)

### Build from source

It is easy to build redis-fs from the source code. It takes four steps:

1. Install `fuse` ([linux](http://fuse.sourceforge.net/), [mac](http://osxfuse.github.io/)). Redis-fs currently works with OS X FUSE 2.8.x but not the 3.x developer preview.
2. Get the redis-fs source code from GitHub

  ```bash
  $ git clone https://github.com/MatthiasWinkelmann/redis-fs.git
  ```

3. Change to the directory with the redis-fs source code and run

  ```bash
  $ make get-deps
  ```

  to install dependencies.

4. Run `make build` and then you can see a binary file in current directory.

### Run Unit Tests

```bash
$ make test
```

## Unmount

### Linux

```bash
$ fusermount -u /tmp/redis
```

### MacOS

```bash
$ diskutil unmount /tmp/redis
```

## License

(The MIT License)
