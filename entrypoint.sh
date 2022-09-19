#!/usr/bin/env bash

env | sed 's/=\(.*\)$/="\1"/g' > env.tmp
set -o allexport
[ -f .env ] && source .env
source env.tmp
set +o allexport
rm env.tmp

exec ./api