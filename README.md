# Hexagonal Saga demo app

hexagonal saga demo app

<img src="/docs/exports/saga-context.png" />
<img src="/docs/exports/overall-service-container.png" />

# Prerequisites

- Docker
- Go 1.19
- MySQL 8.x
- [Wire](https://github.com/google/wire) (for DI)
- [Goose](https://github.com/pressly/goose) (for schema migration)
- [Ginkgo](https://onsi.github.io/ginkgo/), Gomega and [GoMock](https://github.com/golang/mock) for testing
- [swaggo/swag](https://github.com/swaggo/swag)
- [Taskfile](https://taskfile.dev/#/installation)

# Installation

# Run infrastructure

```bash
$ docker-compose --profile backend up -d
```

# Run services

## Copy .env to project root folder

```bash
$ cp env/local.env .env
```

## Schema migration

```bash
$ docker-compose up migrate
```

## Run service

```bash
$ docker-compose up dev
```

## Build swagger docs (Optional)

```bash
$ docker-compose up swagger
```

# Try it out

open swagger on the browser

```bash
$ open http://localhost:8090/swagger/index.html
```

## Create trip record

```bash
$ http --json -v post localhost:8090/v1/trips/ userId:=1 carId:=1 hotelId:=1 flightId:=1

POST /v1/trips/ HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 79
Content-Type: application/json
Host: localhost:8090
User-Agent: HTTPie/2.6.0

{
    "carId": 1,
    "flightId": 1,
    "hotelId": 1,
    "status": "INITIALIZED",
    "userId": 1
}


HTTP/1.1 200 OK
Content-Length: 179
Content-Type: application/json; charset=utf-8
Date: Sun, 28 Aug 2022 12:36:46 GMT

{
    "data": {
        "carId": 1,
        "createdAt": "2022-08-28T21:36:46.825+09:00",
        "flightId": 1,
        "hotelId": 1,
        "id": 1,
        "status": "INITIALIZED",
        "updatedAt": "0001-01-01T00:00:00Z",
        "userId": 1
    },
    "status": true
}
```

## Query created trips

```bash
$ http get localhost:8090/v1/trips/                                                                                                                               dongkyl@DongGyunui-MacBookAir

HTTP/1.1 200 OK
Content-Length: 177
Content-Type: application/json; charset=utf-8
Date: Sun, 28 Aug 2022 12:38:01 GMT

{
    "data": [
        {
            "carId": 1,
            "createdAt": "2022-08-28T19:36:47+07:00",
            "flightId": 1,
            "hotelId": 1,
            "id": 1,
            "status": "INITIALIZED",
            "updatedAt": "0001-01-01T00:00:00Z",
            "userId": 1
        },
        {
            "carId": 1,
            "createdAt": "2022-08-28T19:38:52+07:00",
            "flightId": 1,
            "hotelId": 1,
            "id": 2,
            "status": "INITIALIZED",
            "updatedAt": "0001-01-01T00:00:00Z",
            "userId": 1
        },
        {
            "carId": 1,
            "createdAt": "2022-08-28T19:38:53+07:00",
            "flightId": 1,
            "hotelId": 1,
            "id": 3,
            "status": "INITIALIZED",
            "updatedAt": "0001-01-01T00:00:00Z",
            "userId": 1
        }
    ],
    "status": true
}
```
