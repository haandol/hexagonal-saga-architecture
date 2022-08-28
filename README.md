# HexagonalArchitecture

hexagonal architecture example for gophers

# Prerequisites

- Go 1.18
- MySQL 5.7
- [Wire](https://github.com/google/wire) (for DI)
- [Goose](https://github.com/pressly/goose) (for schema migration)
- [Ginkgo](https://onsi.github.io/ginkgo/), Gomega and [GoMock](https://github.com/golang/mock) for testing

# Installation

```bash
$ go mod tidy && go mod vendor
```

# Run infrastructure

```bash
$ docker-compose --profile backend up -d
```

# Schema migration

Install Goose

```bash
$ go install github.com/pressly/goose/v3/cmd/goose@latest
```

```bash
$ ./scripts/migrate.sh
```

# Run Swagger

install [swaggo/swag](https://github.com/swaggo/swag) 

```bash
$ go install github.com/swaggo/swag/cmd/swag@1.8.4
```

```bash
$ ./script/swagger.sh
```

# Run server

run order service

```bash
$ cd cmd/order
$ go run main.go

2022-08-28T08:35:14.686+0900	INFO	order/main.go:55
```

run server and open swagger on the browser

```bash
$ open http://localhost:8080/swagger/index.html
```

## Create trip record

```bash
$ http --json -v post localhost:8080/v1/trips/ userId:=1 carId:=1 hotelId:=1 flightId:=1 status=INITIALIZED

POST /v1/trips/ HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 79
Content-Type: application/json
Host: localhost:8080
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
$ http get localhost:8080/v1/trips/                                                                                                                               dongkyl@DongGyunui-MacBookAir

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