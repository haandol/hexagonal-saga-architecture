#!/bin/sh

set -o allexport
[ -f ./.env ] && . ./.env
set +o allexport

cd init

echo $DB_USERNAME
echo $DB_HOST:$DB_PORT
export GOOSE_DBSTRING="$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true"

echo $(goose mysql status)

exec goose mysql $1
