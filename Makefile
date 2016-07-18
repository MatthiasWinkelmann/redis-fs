APPNAME=redis-fs
SERVICEPATH=github.com/MatthiasWinkelmann

ENTRYPOINT=main.go
GOPATH=$(shell pwd)

define build
	echo $(APPNAME)-$(1)-$(2); \
	GO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o "bin/$(APPNAME)-$(1)-$(2)" $(ENTRYPOINT);
endef

RUN_ARGS=$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

test:
	@go test -v ./redisfs

run:
	@go run $(ENTRYPOINT) $(RUN_ARGS)

cross:
	@$(call build,linux,amd64)
	@$(call build,linux,386)
	@$(call build,linux,arm)
	@$(call build,darwin,amd64)

build:
	@go build

get-deps:
	@go get github.com/codegangsta/cli
	@go get github.com/hanwen/go-fuse/fuse
	@go get github.com/visionmedia/go-debug
	@go get github.com/garyburd/redigo/redis
	@go get github.com/smartystreets/goconvey/convey
	@mkdir -p src/$(SERVICEPATH)
	ln -s $(GOPATH) src/$(SERVICEPATH)/$(APPNAME)

clean:
	-@rm -rf bin src pkg
