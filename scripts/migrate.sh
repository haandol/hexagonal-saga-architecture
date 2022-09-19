#!/usr/bin/env sh

set -o allexport
[ -f .env ] && source .env
set +o allexport

cd init

echo $DB_USERNAME
echo $DB_HOST:$DB_PORT
export GOOSE_DBSTRING="host=$DB_HOST user=$DB_USERNAME password=$DB_PASSWORD dbname=$DB_NAME port=$DB_PORT sslmode=disable"

echo $(goose postgres status)

exec goose postgres up
