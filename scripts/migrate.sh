#!/bin/sh

set -o allexport
[ -f ./.env ] && . ./.env
set +o allexport

echo $DB_USERNAME
echo $DB_HOST:$DB_PORT
export GOOSE_DBSTRING="$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true"

echo $(goose -dir init mysql status)

exec goose -dir init mysql $1
