#!/bin/sh

set -o allexport
[ -f ./.env ] && . ./.env
set +o allexport

if [ ! -z $DB_SECRET_ID ]; then
  echo 'Getting DB credentials from AWS Secrets Manager'
  export DB_USERNAME=$(aws secretsmanager get-secret-value --secret-id $DB_SECRET_ID --query SecretString | jq -r | jq -r '.username')
  export DB_PASSWORD=$(aws secretsmanager get-secret-value --secret-id $DB_SECRET_ID --query SecretString | jq -r | jq -r '.password')
fi

echo $DB_USERNAME
echo $DB_HOST:$DB_PORT
export GOOSE_DBSTRING="$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true"

echo $(goose -dir init mysql status)

exec goose -dir init mysql $1
