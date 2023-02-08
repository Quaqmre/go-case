# Getir Case

## Requirements

* Go (1.16)
* Docker
* Docker Compose
* GNU make (optional)

## Build & Run

With Docker Compose

```shell
cd <project dir>

# build
docker-compose -f docker-compose.yml build

# run
docker-compose -f docker-compose.yml up
# or run in background
docker-compose -f docker-compose.yml up -d

# stop
docker-compose -f docker-compose.yml down
```

With make

```shell
cd <project dir>

make build
make debug # run
make run # run in background
make stop
```

P.S. Tested in Arch Linux

### Testing

Without make

```shell
go mod download # only once to download modules
ginkgo ./...
```

With make

```shell
make dep # only once to download modules
make test
```
