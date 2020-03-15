FROM golang:1.14 as builder

RUN mkdir src/redis-fs
WORKDIR src/redis-fs

COPY Makefile .
COPY go.mod .
RUN make get-deps

COPY . .
RUN make build
RUN pwd
RUN ls -l redis-fs

# Can't run this without a running redis instance
# RUN make test


FROM ubuntu:bionic

COPY --from=builder /go/src/redis-fs/redis-fs /usr/bin/

ENTRYPOINT /usr/bin/redis-fs
