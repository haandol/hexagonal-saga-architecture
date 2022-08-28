#!/usr/bin/env sh

set -o allexport
[ -f env/local.env ] && source env/local.env
set +o allexport

cd migration

echo $DB_USERNAME
echo $DB_HOST:$DB_PORT
export GOOSE_DBSTRING="$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true"

echo $(goose mysql status)

exec goose mysql up