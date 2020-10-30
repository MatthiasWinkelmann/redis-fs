redis-fs
===========

Note: This is largely abandoned by me, as it was by others, before. So it probably has issues. 

The main reason is that it's slow. I couldn't get even close to rotating disk HDD transfer rates.  believe that was the whole point, back when I still knew what I was trying to do...

## Sawing off the branch I'm sitting on...

Check out the forks ion the Insights tab, above, and you might find someone who took this idea and ran with it. [LexVocoder](https://github.com/LexVocoder) and [Promaethius](https://github.com/Promaethius) were the last being frudtrated by my inactivity.

redis-fs lets you mount a Redis database as a filesystem. It is based on redis-mount by Po-Ying Chen which was deleted from Github for unknown reasons.


## Usage

```bash
redis-fs 0.2.0
$ redis-fs ~/redis

--host, -h   localhost    Redis host name
--port, -p   6379         Redis port number
--auth, -a                Redis password
--sep, -s    :            Redis key separator
```

## Installation

### Download binary file

* [mac-amd64](https://github.com/MatthiasWinkelmann/redis-fs/raw/master/releases/download/redis-fs-darwin-amd64)
* [linux-amd64](https://github.com/MatthiasWinkelmann/redis-fs/raw/master/releases/download/redis-fs-linux-amd64)
* [linux-386](https://github.com/MatthiasWinkelmann/redis-fs/raw/master/releases/download/redis-fs-linux-386)
* [linux-arm](https://github.com/MatthiasWinkelmann/redis-fs/raw/master/releases/download/redis-fs-linux-arm)

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
