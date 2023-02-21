#!/bin/sh

set -o allexport
[ -f ./.env ] && . ./.env
set +o allexport

if [ ! -z $DB_SECRET_ID ]; then
  echo "Getting DB credentials from AWS Secrets Manager with id - $DB_SECRET_ID"
  aws secretsmanager get-secret-value --secret-id $DB_SECRET_ID --query SecretString | jq -r > /tmp/db-secret.json
  export DB_HOST=$(cat /tmp/db-secret.json | jq -r '.host')
  export DB_PORT=$(cat /tmp/db-secret.json | jq -r '.port')
  export DB_NAME=$(cat /tmp/db-secret.json | jq -r '.dbname')
  export DB_USERNAME=$(cat /tmp/db-secret.json | jq -r '.username')
  export DB_PASSWORD=$(cat /tmp/db-secret.json | jq -r '.password')
  rm /tmp/db-secret.json
fi

echo $DB_USERNAME
echo $DB_HOST:$DB_PORT
export GOOSE_DBSTRING="$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true"

echo $(goose -dir init mysql status)

exec goose -dir init mysql $1
