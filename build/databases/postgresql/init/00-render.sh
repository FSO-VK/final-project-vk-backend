#!/bin/sh
set -e

# Писать во временную папку, а не в смонтированную
envsubst < /docker-entrypoint-initdb.d/01-init.sql.template \
  > /tmp/01-init.sql

# Выполнить SQL из временного файла
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" \
  -f /tmp/01-init.sql