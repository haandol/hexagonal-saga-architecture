# Hexagonal Saga demo app

hexagonal saga demo app

<img src="/docs/exports/saga-context.png" />
<img src="/docs/exports/overall-service-container.png" />

# Prerequisites

- Docker
- Go 1.23+
- MySQL 8.x
- [Wire](https://github.com/google/wire) (for DI)
- [Goose](https://github.com/pressly/goose) (for schema migration)
- [Ginkgo](https://onsi.github.io/ginkgo/), Gomega and [GoMock](https://github.com/golang/mock) for testing
- [swaggo/swag](https://github.com/swaggo/swag)
- [Taskfile](https://taskfile.dev/#/installation)

# Installation

# Run infrastructure

```bash
docker-compose --profile backend up -d
```

# Run services

## Copy .env to project root folder

```bash
cp env/local.env .env
```

## Schema migration

```bash
docker-compose up migrate
```

## Run service

```bash
docker-compose up dev
```

## Build swagger docs (Optional)

```bash
docker-compose up swagger
```

# Try it out

open swagger on the browser

```bash
open http://localhost:8090/swagger/index.html
```

## Create trip record

```bash
$ http --json -v post localhost:8090/v1/trips/ < payload/create-trip.json

POST /v1/trips/ HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 76
Content-Type: application/json
Host: localhost:8090
User-Agent: HTTPie/3.2.2

{
    "carId": 1,
    "flightId": 1,
    "hotelId": 1,
    "id": 1,
    "userId": 1
}


HTTP/1.1 200 OK
Content-Length: 229
Content-Type: application/json; charset=utf-8
Date: Sat, 09 Sep 2023 09:17:13 GMT

{
    "data": {
        "carBookingId": 0,
        "carId": 1,
        "createdAt": "2023-09-09T09:17:13.87Z",
        "flightBookingId": 0,
        "flightId": 1,
        "hotelBookingId": 0,
        "hotelId": 1,
        "id": 5,
        "status": "Initialized",
        "updatedAt": "0001-01-01T00:00:00Z",
        "userId": 1
    },
    "status": true
}
```

## Query created trips

```bash
$ http get localhost:8090/v1/trips/

HTTP/1.1 200 OK
Content-Length: 1044
Content-Type: application/json; charset=utf-8
Date: Sat, 09 Sep 2023 09:18:01 GMT

{
    "data": [
        {
            "carBookingId": 2,
            "carId": 1,
            "createdAt": "2023-09-09T18:13:43+09:00",
            "flightBookingId": 2,
            "flightId": 1,
            "hotelBookingId": 2,
            "hotelId": 1,
            "id": 2,
            "status": "Booked",
            "updatedAt": "0001-01-01T00:00:00Z",
            "userId": 1
        },
        {
            "carBookingId": 1,
            "carId": 1,
            "createdAt": "2023-09-09T18:13:36+09:00",
            "flightBookingId": 1,
            "flightId": 1,
            "hotelBookingId": 1,
            "hotelId": 1,
            "id": 1,
            "status": "Booked",
            "updatedAt": "0001-01-01T00:00:00Z",
            "userId": 1
        }
    ],
    "status": true
}
```
