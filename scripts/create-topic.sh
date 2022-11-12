#!/bin/sh

set -o allexport
[ -f ./.env ] && . ./.env
set +o allexport

docker run --rm haandol/kafka-cli:2.6.2 kafka-topics.sh --create --topic trip-service --partitions 10 --replication-factor 3 --bootstrap-server $KAFKA_SEEDS || exit 1
docker run --rm haandol/kafka-cli:2.6.2 kafka-topics.sh --create --topic saga-service --partitions 10 --replication-factor 3 --bootstrap-server $KAFKA_SEEDS || exit 1
docker run --rm haandol/kafka-cli:2.6.2 kafka-topics.sh --create --topic car-service --partitions 10 --replication-factor 3 --bootstrap-server $KAFKA_SEEDS || exit 1
docker run --rm haandol/kafka-cli:2.6.2 kafka-topics.sh --create --topic flight-service --partitions 10 --replication-factor 3 --bootstrap-server $KAFKA_SEEDS || exit 1
docker run --rm haandol/kafka-cli:2.6.2 kafka-topics.sh --create --topic hotel-service --partitions 10 --replication-factor 3 --bootstrap-server $KAFKA_SEEDS || exit 1
