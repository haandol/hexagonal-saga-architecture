FROM golang:1.20 AS builder

WORKDIR /src

# set environment path
ENV PATH /go/bin:$PATH
ENV GONOSUMDB github.com/haandol
ENV GOPRIVATE github.com/haandol

# manage dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . ./

ARG BUILD_TAG
ARG APP_NAME
ARG TARGETOS=linux
ARG TARGETARCH=arm64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-X main.BuildTag=$BUILD_TAG -s -w" -o /go/bin/app ./cmd/${APP_NAME}

FROM alpine:3.17.1 AS server
ARG GIT_COMMIT=undefined
LABEL git_commit=$GIT_COMMIT

RUN apk --no-cache add curl
RUN apk --no-cache add tzdata
RUN echo "Asia/Seoul" >  /etc/timezone

WORKDIR /
COPY --chown=0:0 --from=builder /go/bin/app /
COPY --chown=0:0 --from=builder /src/docker-entrypoint.sh /
COPY --chown=0:0 --from=builder /src/env/dev.env /.env

ARG APP_PORT
EXPOSE ${APP_PORT}

ENTRYPOINT ["/docker-entrypoint.sh"]

# for schema migration
FROM haandol/goose:3.9.0 as migration
ARG GIT_COMMIT=undefined
LABEL git_commit=$GIT_COMMIT

RUN apk add --no-cache aws-cli jq

WORKDIR /
COPY ./init /init
COPY ./env/dev.env /.env
COPY ./scripts/migrate.sh /migrate.sh

CMD ["/migrate.sh", "up"]